package repository

import (
	"context"

	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/libraries/result"
)

type SettingRepository interface {
	List(ctx context.Context) result.Result
	Get(ctx context.Context, key string) result.Result
	Update(ctx context.Context, d *model.Setting) result.Result
}
