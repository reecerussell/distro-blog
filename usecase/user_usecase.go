package usecase

import (
	"context"
	"net/http"

	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/domain/service"
	"github.com/reecerussell/distro-blog/libraries/normalization"
	"github.com/reecerussell/distro-blog/libraries/result"
	"github.com/reecerussell/distro-blog/password"
)

// UserUsecase is a high-level interface used to manage the user domain.
type UserUsecase interface {
	Create(ctx context.Context, cu *dto.CreateUser) result.Result
}

// userUsecase is an implementation of the UserUsecase interface.
type userUsecase struct {
	repo    repository.UserRepository
	serv    *service.UserService
	norm    normalization.Normalizer
	pwdServ password.Service
}

// NewUserUsecase returns a new instance of the UserUsecase
// interface with the given repository.
func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	serv := service.NewUserService(repo)
	norm := normalization.New()
	pwdServ := password.New()

	return &userUsecase{
		repo:    repo,
		serv:    serv,
		norm:    norm,
		pwdServ: pwdServ,
	}
}

// Create creates a new user domain record, ensuring the data is valid.
func (u *userUsecase) Create(ctx context.Context, cu *dto.CreateUser) result.Result {
	usr, err := model.NewUser(cu, u.pwdServ, u.norm)
	if err != nil {
		return result.Failure(err).WithStatusCode(http.StatusBadRequest)
	}

	unique, _, _, err := u.serv.EnsureEmailIsUnique(ctx, usr).Deconstruct()
	if !unique {
		return result.Failure(err).WithStatusCode(http.StatusBadRequest)
	}

	success, _, _, err := u.repo.Add(ctx, usr).Deconstruct()
	if !success {
		return result.Failure(err).WithStatusCode(http.StatusInternalServerError)
	}

	return result.Ok()
}
