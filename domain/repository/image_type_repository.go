package repository

import (
	"context"

	"github.com/reecerussell/distro-blog/libraries/result"
)

type ImageTypeRepository interface {
	GetByName(ctx context.Context, name string) result.Result
}

