package repository

import (
	"context"

	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/libraries/result"
)

type ImageRepository interface {
	Get(ctx context.Context, id string) result.Result
	Add(ctx context.Context, i *model.Image) result.Result
	Delete(ctx context.Context, id string) result.Result
}
