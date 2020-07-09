package usecase

import (
	"context"
	"fmt"
	"net/http"
	"os"

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

	err = u.stg.Set(it.GetID(), data)
	if err != nil {
		return result.Failure(err)
	}

	res := u.ir.Add(ctx, img)
	if !res.IsOk() {
		return res
	}

	return result.Ok().WithValue(img)
}