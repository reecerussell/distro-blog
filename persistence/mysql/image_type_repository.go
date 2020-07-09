package mysql

import (
	"context"
	"fmt"
	"net/http"

	"github.com/reecerussell/distro-blog/domain/datamodel"
	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/normalization"
	"github.com/reecerussell/distro-blog/libraries/result"
)

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
	name = r.norm.Normalize(name)

	dm, err := r.db.Read(ctx, query, imageTypeReader, name)
	if err != nil {
		msg := fmt.Sprintf("no image type exists with name '%s'", name)
		return result.Failure(msg).WithStatusCode(http.StatusNotFound)
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
