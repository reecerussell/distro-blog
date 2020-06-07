package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/result"
)

// UserService is used to provide the user domain with extra functionality,
// such as validation methods which don't belong to the domain layer.
type UserService struct {
	repo repository.UserRepository
}

// NewUserService returns a new instance of UserService with the given repository.
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// EnsureEmailIsUnique ensures the given user's email address is unique by
// ensuring it has not be previously used in the database.
func (s *UserService) EnsureEmailIsUnique(ctx context.Context, u *model.User) result.Result {
	success, _, count, err := s.repo.CountByEmail(ctx, u).Deconstruct()
	if !success {
		return result.Failure(err)
	}

	fmt.Printf("EnsureEmailIsUnique: %d\n", count)

	if count.(int64) > 0 {
		msg := fmt.Sprintf("The email address '%s' has already been taken.", u.Email())
		return result.Failure(msg).WithStatusCode(http.StatusBadRequest)
	}

	return result.Ok()
}
