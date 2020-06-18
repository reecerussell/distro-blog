package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/reecerussell/distro-blog/auth"
	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/result"
	"github.com/reecerussell/distro-blog/password"
)

type AuthUsecase interface {
	Token(ctx context.Context, cred *dto.UserCredential) result.Result
	Verify(ctx context.Context, tokenData []byte) result.Result
}

type authUsecase struct {
	repo repository.UserRepository
	pwd password.Service
	auth *auth.Service
}

func NewAuthUsecase(repo repository.UserRepository) AuthUsecase {
	return &authUsecase{
		repo: repo,
		pwd: password.New(),
		auth: auth.New(),
	}
}

func (u *authUsecase) Token(ctx context.Context, cred *dto.UserCredential) result.Result {
	defaultErr := "Email and/or password is incorrect."

	success, status, value, err := u.repo.GetByEmail(ctx, cred.Email).Deconstruct()
	if !success {
		if status == http.StatusNotFound {
			return result.Failure(defaultErr).WithStatusCode(http.StatusBadRequest)
		}

		return result.Failure(err).WithStatusCode(status)
	}

	user := value.(*model.User)
	err = user.VerifyPassword(cred.Password, u.pwd)
	if err != nil {
		return result.Failure(defaultErr).WithStatusCode(http.StatusBadRequest)
	}

	scopes := user.Scopes()
	scopeNames := make([]string, len(scopes))

	for i, s := range scopes {
		scopeNames[i] = s.Name()
	}

	claims := map[string]interface{}{
		auth.ClaimTypeEmail: user.NormalizedEmail(),
		auth.ClaimTypeUserId: user.ID(),
		auth.ClaimTypeScopes: scopeNames,
	}

	now := time.Now().UTC()
	exp := now.Add(time.Second * 3600)

	t := u.auth.NewToken(ctx).
		SetNotBefore(now).
		SetIssuedAt(now).
		SetExpiry(exp).
		AddClaims(claims).
		Build()

	ac := auth.NewAccessToken(t, exp)
	return result.Ok().WithValue(ac)
}

func (u *authUsecase) Verify(ctx context.Context, tokenData []byte) result.Result {
	ok := u.auth.VerifyToken(ctx, tokenData)
	if !ok {
		return result.Failure("invalid token").
			WithStatusCode(http.StatusUnauthorized)
	}

	return result.Ok()
}