package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/reecerussell/distro-blog/domain/datamodel"
	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/logging"
	"github.com/reecerussell/distro-blog/libraries/normalization"
	"github.com/reecerussell/distro-blog/libraries/result"
)

var errMsgImageTypeServerError = "IMAGE_TYPE_SERVER_ERROR"

type imageTypeRepository struct {
	db *database.MySQL
	norm normalization.Normalizer
}

func NewImageTypeRepository(db *database.MySQL) repository.ImageTypeRepository {
	return &imageTypeRepository{
		db: db,
		norm: normalization.New(),
	}
}

func (r *imageTypeRepository) GetByName(ctx context.Context, name string) result.Result{
	const query string = "CALL `get_image_type_by_name`(?);"
	dm, err := r.db.Read(ctx, query, imageTypeReader, r.norm.Normalize(name))
	if err != nil {
		if err == sql.ErrNoRows {
			msg := fmt.Sprintf("no image type exists with name '%s'", name)
			return result.Failure(msg).WithStatusCode(http.StatusNotFound)
		}

		logging.Error(err)
		return result.Failure(errMsgImageTypeServerError)
	}

	it := model.ImageTypeFromDataModel(dm.(*datamodel.ImageType))
	return result.Ok().WithValue(it)
}

func imageTypeReader(s database.ScannerFunc) (interface{}, error) {
	var dm datamodel.ImageType
	err := s(
		&dm.ID,
		&dm.Name,
	)
	if err != nil {
		return nil, err
	}

	return &dm, nil
}
