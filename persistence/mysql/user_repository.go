package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/reecerussell/distro-blog/domain/datamodel"
	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/result"
	"net/http"
)

type userRepository struct {
	db *database.MySQL
}

func NewUserRepository(db *database.MySQL) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) List(ctx context.Context) result.Result {
	query := "SELECT * FROM `view_user_list`;"
	items, err := r.db.Multiple(ctx, query, userReader)
	if err != nil {
		return result.Failure(err)
	}

	dtos := make([]*dto.UserListItem, len(items))

	for i, item := range items {
		dto := item.(dto.UserListItem)
		dtos[i] = &dto
	}

	return result.Ok().WithValue(dtos)
}

func userReader(s database.ScannerFunc) (interface{}, error) {
	var dto dto.UserListItem
	err := s(
		&dto.ID,
		&dto.Name,
		&dto.Email,
	)
	if err != nil {
		return nil, err
	}

	return dto, nil
}

func (r *userRepository) Get(ctx context.Context, id string) result.Result {
	const query string = "CALL `get_user`(?);"

	dm, err := r.db.Read(ctx, query, func(s database.ScannerFunc) (interface{}, error) {
		var dm datamodel.User
		err := s(
			&dm.ID,
			&dm.Firstname,
			&dm.Lastname,
			&dm.Email,
			&dm.NormalizedEmail,
			&dm.PasswordHash,
		)
		if err != nil {
			return nil, err
		}

		return &dm, nil
	}, id)
	if err != nil && err != sql.ErrNoRows{
		return result.Failure(err)
	}

	if dm == nil || err == sql.ErrNoRows {
		msg := fmt.Sprintf("No user was found with id '%s'", id)
		return result.Failure(msg).WithStatusCode(http.StatusNotFound)
	}

	u := model.UserFromDataModel(dm.(*datamodel.User))

	return result.Ok().
		WithValue(u)
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

	_, err := r.db.Execute(ctx, query, args...)
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

// Update modifies an existing user record in the database, with the updated domain model.
func (r *userRepository) Update(ctx context.Context, u *model.User) result.Result {
	const query string = "CALL `update_user`(?,?,?,?,?);"
	dm := u.DataModel()
	args := []interface{}{
		dm.ID,
		dm.Firstname,
		dm.Lastname,
		dm.Email,
		dm.NormalizedEmail,
	}

	ra, err := r.db.Execute(ctx, query, args...)
	if err != nil {
		return result.Failure(err)
	}

	if ra < 1 {
		msg := fmt.Sprintf("No user with id '%s' was updated.", dm.ID)
		return result.Failure(msg).WithStatusCode(http.StatusNotFound)
	}

	return result.Ok()
}