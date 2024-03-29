package usecase

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/reecerussell/distro-blog/libraries/logging"
	"net/http"
	"os"
	"strings"

	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/result"
	"github.com/reecerussell/distro-blog/libraries/storage"
)

var mediaS3BucketName = os.Getenv("MEDIA_BUCKET_NAME")

type MediaUsecase interface {
	// Upload puts the data into storage, and returns a result
	// with the value set as the id of the newly inserted database record.
	Upload(ctx context.Context, data []byte) result.Result

	Delete(ctx context.Context, id string) result.Result
	DownloadForLambda(ctx context.Context, id string) events.APIGatewayProxyResponse
}

func NewMediaUsecase(ir repository.ImageRepository, itr repository.ImageTypeRepository) (MediaUsecase, error) {
	stg, err := storage.New(mediaS3BucketName)
	if err != nil {
		return nil, err
	}

	return &mediaUsecase{
		ir:  ir,
		itr: itr,
		stg: stg,
	}, nil
}

type mediaUsecase struct {
	ir repository.ImageRepository
	itr repository.ImageTypeRepository
	stg *storage.Service
}

func (u *mediaUsecase) Upload(ctx context.Context, data []byte) result.Result {
	mt := http.DetectContentType(data)

	success, status, value, err := u.itr.GetByName(ctx, mt).Deconstruct()
	if !success {
		if status == http.StatusNotFound {
			msg := fmt.Sprintf("Image mime type '%s' is unsupported.", mt)
			return result.Failure(msg).WithStatusCode(http.StatusBadRequest)
		}

		return result.Failure(err).WithStatusCode(status)
	}

	it := value.(*model.ImageType)
	img := model.NewImage(it)

	err = u.stg.SetImage(img.GetID(), mt, data)
	if err != nil {
		return result.Failure(err)
	}

	res := u.ir.Add(ctx, img)
	if !res.IsOk() {
		return res
	}

	return result.Ok().WithValue(img)
}

func (u *mediaUsecase) Delete(ctx context.Context, id string) result.Result {
	res := u.ir.Delete(ctx, id)
	if !res.IsOk() {
		return res
	}

	err := u.stg.Delete(id)
	if err != nil {
		return result.Failure(err)
	}

	return result.Ok()
}

func (u *mediaUsecase) DownloadForLambda(ctx context.Context, id string) events.APIGatewayProxyResponse {
	success, status, imageType, err := u.ir.GetType(ctx, id).Deconstruct()
	if !success {
		logging.Errorf("Failed to get image type: %v\n", err)
		return events.APIGatewayProxyResponse{StatusCode: status}
	}

	data, err := u.stg.Get(id)
	if err != nil {
		logging.Errorf("Failed to download image: %v\n", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": strings.ToLower(imageType.(string)),
			"Cache-Control": "max-age=604800",
		},
		StatusCode: http.StatusOK,
		Body: base64.StdEncoding.EncodeToString(data),
		IsBase64Encoded: true,
	}
}