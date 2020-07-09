package repository

import (
	"context"

	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/libraries/result"
)

type ImageRepository interface {
	Add(ctx context.Context, i *model.Image) result.Result
}
