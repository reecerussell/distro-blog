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
	List(ctx context.Context) result.Result
	Get(ctx context.Context, id string) result.Result
	Create(ctx context.Context, cu *dto.CreateUser) result.Result
	Update(ctx context.Context, uu *dto.UpdateUser) result.Result
	Token(ctx context.Context, cred *dto.UserCredential) result.Result
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

// List returns a result containing a list of users, from the repository.
func (u *userUsecase) List(ctx context.Context) result.Result {
	return u.repo.List(ctx)
}

// Get retrieves a specific user from the database and returns a result with a dto value.
func (u *userUsecase) Get(ctx context.Context, id string) result.Result {
	res := u.repo.Get(ctx, id)
	if !res.IsOk(){
		return res
	}

	_, _, value, _ := res.Deconstruct()
	user := value.(*model.User)

	return result.Ok().WithValue(user.DTO())
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

// Update updates a specific user record, by retrieving the record from the
// database, then updating the domain model and database record.
func (u *userUsecase) Update(ctx context.Context, uu *dto.UpdateUser) result.Result {
	success, status, value, err := u.repo.Get(ctx, uu.ID).Deconstruct()
	if !success{
		return result.Failure(err).WithStatusCode(status)
	}

	user := value.(*model.User)
	err = user.Update(uu, u.norm)
	if err != nil {
		return result.Failure(err).WithStatusCode(http.StatusBadRequest)
	}

	success, status, _, err = u.repo.Update(ctx, user).Deconstruct()
	if !success {
		return result.Failure(err).WithStatusCode(status)
	}

	return result.Ok()
}

func (u *userUsecase) Token(ctx context.Context, cred *dto.UserCredential) result.Result {

}