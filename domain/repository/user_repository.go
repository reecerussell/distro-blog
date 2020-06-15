package repository

import (
	"context"

	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/libraries/result"
)

// UserRepository is a high-level interface used to manage
// user persistence and handle interaction with a data source.
type UserRepository interface {
	List(ctx context.Context) result.Result
	Get(ctx context.Context, id string) result.Result
	Add(ctx context.Context, u *model.User) result.Result
	CountByEmail(ctx context.Context, u *model.User) result.Result
	Update(ctx context.Context, u *model.User) result.Result
	Delete(ctx context.Context, id string) result.Result
}
