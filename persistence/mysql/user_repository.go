package mysql

import (
	"context"

	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/result"
)

type userRepository struct {
	db *database.MySQL
}

func NewUserRepository(db *database.MySQL) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Add(ctx context.Context, u *model.User) result.Result {
	dm := u.DataModel()
	query := "CALL `create_user`(?,?,?,?,?,?);"
	args := []interface{}{
		dm.ID,
		dm.Firstname,
		dm.Lastname,
		dm.Email,
		dm.NormalizedEmail,
		dm.PasswordHash,
	}

	err := r.db.Execute(ctx, query, args...)
	if err != nil {
		return result.Failure(err)
	}

	return result.Ok()
}

// CountByEmail counts the number of users in the database with the given email. If successful,
// an "Ok" result will be returned with an int64 value.
func (r *userRepository) CountByEmail(ctx context.Context, u *model.User) result.Result {
	const query string = "CALL `count_users_by_email`(?, ?);"

	dm := u.DataModel()
	c, err := r.db.Count(ctx, query, dm.NormalizedEmail, dm.ID)
	if err != nil {
		return result.Failure(err)
	}

	return result.Ok().WithValue(c)
}
