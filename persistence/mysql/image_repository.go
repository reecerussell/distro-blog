package mysql

import (
	"context"

	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/result"
)

type imageRepository struct {
	db *database.MySQL
}

func NewImageRepository(db *database.MySQL) repository.ImageRepository {
	return &imageRepository{
		db: db,
	}
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
		return result.Failure(err)
	}

	return result.Ok()
}
