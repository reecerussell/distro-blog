package mysql

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/reecerussell/distro-blog/domain/datamodel"
	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/logging"
	"github.com/reecerussell/distro-blog/libraries/result"
)

var (
	errMsgImageNotFound = "IMAGE_NOT_FOUND"
	errMsgImageDbError = "IMAGE_SERVER_ERROR"
)

type imageRepository struct {
	db *database.MySQL
}

func NewImageRepository(db *database.MySQL) repository.ImageRepository {
	return &imageRepository{
		db: db,
	}
}

func (r *imageRepository) Get(ctx context.Context, id string) result.Result {
	const query string = "CALL `get_image`(?);"
	dm, err := r.db.Read(ctx, query, imageReader, id)
	if err != nil && err != sql.ErrNoRows {
		logging.Error(err)
		return result.Failure(errMsgImageDbError)
	}

	if dm == nil || err == sql.ErrNoRows {
		return result.Failure(errMsgImageNotFound).WithStatusCode(http.StatusNotFound)
	}

	img := model.ImageFromDataModel(dm.(*datamodel.Image))
	return result.Ok().WithValue(img)
}

func imageReader(s database.ScannerFunc) (interface{}, error) {
	var dm datamodel.Image
	err := s(
		&dm.ID,
		&dm.TypeID,
		&dm.AlternativeText,
	)
	if err != nil {
		return nil, err
	}

	return &dm, nil
}

func (r *imageRepository) Add(ctx context.Context, i *model.Image) result.Result {
	const query string = "CALL `create_image`(?,?);"
	dm := i.DataModel()
	args := []interface{}{
		dm.ID,
		dm.TypeID,
	}

	_, err := r.db.Execute(ctx, query, args...)
	if err != nil {
		logging.Error(err)
		return result.Failure(errMsgImageDbError)
	}

	return result.Ok()
}

func (r *imageRepository) Delete(ctx context.Context, id string) result.Result {
	const query string = "CALL `delete_image`(?);"
	ra, err := r.db.Execute(ctx, query, id)
	if err != nil {
		logging.Error(err)
		return result.Failure(errMsgImageDbError)
	}

	if ra < 1 {
		return result.Failure(errMsgImageNotFound).WithStatusCode(http.StatusNotFound)
	}

	return result.Ok()
}

func (r *imageRepository) GetType(ctx context.Context, id string) result.Result {
	const query string = "CALL `get_image_type`(?);"
	imageType, err := r.db.Read(ctx, query, getTypeReader, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return result.Failure(errMsgImageNotFound).WithStatusCode(http.StatusNotFound)
		}

		logging.Error(err)
		return result.Failure(errMsgImageDbError)
	}

	return result.Ok().WithValue(imageType)
}

func getTypeReader(s database.ScannerFunc) (interface{}, error) {
	var imageType string
	err := s(&imageType)
	return imageType, err
}