package repository

import (
	"context"

	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/libraries/result"
)

// UserRepository is a high-level interface used to manage
// user persistance and handle interaction with a data source.
type UserRepository interface {
	Add(ctx context.Context, u *model.User) result.Result
	CountByEmail(ctx context.Context, u *model.User) result.Result
}
