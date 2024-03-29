package usecase

import (
	"context"
	"github.com/reecerussell/distro-blog/libraries/contextkey"
	"github.com/reecerussell/distro-blog/libraries/logging"
	"net/http"
	"strings"

	"github.com/reecerussell/distro-blog/auth"
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
	Get(ctx context.Context, id string, expand ...string) result.Result
	Create(ctx context.Context, cu *dto.CreateUser) result.Result
	Update(ctx context.Context, uu *dto.UpdateUser) result.Result
	Delete(ctx context.Context, id string) result.Result
	ChangePassword(ctx context.Context, d *dto.ChangePassword) result.Result
	ResetPassword(ctx context.Context, id string) result.Result
}

// userUsecase is an implementation of the UserUsecase interface.
type userUsecase struct {
	repo    repository.UserRepository
	serv    *service.UserService
	norm    normalization.Normalizer
	pwdServ password.Service
	auth *auth.Service
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
func (u *userUsecase) Get(ctx context.Context, id string, expand ...string) result.Result {
	res := u.repo.Get(ctx, id)
	if !res.IsOk(){
		return res
	}

	_, _, value, _ := res.Deconstruct()
	user := value.(*model.User).DTO()

	for _, e := range expand {
		switch strings.ToLower(e) {
		case "audit":
			logging.Debugf("Expanded Audit.\n")
			success, _, audit, err := u.repo.GetAudit(ctx, id).Deconstruct()
			if success {
				user.Audit = audit.([]*dto.UserAudit)
			} else {
				logging.Errorf("An error occurred while getting the user's audit data: %v", err)
			}
		}
	}

	return result.Ok().WithValue(user)
}

// Create creates a new user domain record, ensuring the data is valid.
func (u *userUsecase) Create(ctx context.Context, cu *dto.CreateUser) result.Result {
	usr, err := model.NewUser(ctx, cu, u.pwdServ, u.norm)
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
	err = user.Update(ctx, uu, u.norm)
	if err != nil {
		return result.Failure(err).WithStatusCode(http.StatusBadRequest)
	}

	success, status, _, err = u.repo.Update(ctx, user).Deconstruct()
	if !success {
		return result.Failure(err).WithStatusCode(status)
	}

	return result.Ok()
}

// Delete deletes the user from the current database.
func (u *userUsecase) Delete(ctx context.Context, id string) result.Result {
	// ensure exists
	success, status, _, err := u.repo.Get(ctx, id).Deconstruct()
	if !success {
		return result.Failure(err).WithStatusCode(status)
	}

	return u.repo.Delete(ctx, id)
}

// ChangePassword is used to change the current user's password.
func (u *userUsecase) ChangePassword(ctx context.Context, d *dto.ChangePassword) result.Result {
	uid := ctx.Value(contextkey.ContextKey("user_id"))
	if uid == nil {
		return result.Failure("User must be logged in to change password.").
			WithStatusCode(http.StatusUnauthorized)
	}

	success, status, value, err := u.repo.Get(ctx, uid.(string)).Deconstruct()
	if !success{
		return result.Failure(err).WithStatusCode(status)
	}

	user := value.(*model.User)
	success, status, _, err = user.ChangePassword(ctx, d, u.pwdServ).Deconstruct()
	if !success{
		return result.Failure(err).WithStatusCode(status)
	}

	success, status, _, err = u.repo.Update(ctx, user).Deconstruct()
	if !success {
		return result.Failure(err).WithStatusCode(status)
	}

	return result.Ok()
}

// ResetPassword resets the password for the user with the given id.
func (u *userUsecase) ResetPassword(ctx context.Context, id string) result.Result {
	success, status, value, err := u.repo.Get(ctx, id).Deconstruct()
	if !success{
		return result.Failure(err).WithStatusCode(status)
	}

	user := value.(*model.User)
	pwd := user.ResetPassword(ctx, u.pwdServ)

	success, status, _, err = u.repo.Update(ctx, user).Deconstruct()
	if !success {
		return result.Failure(err).WithStatusCode(status)
	}

	return result.Ok().WithValue(pwd)
}