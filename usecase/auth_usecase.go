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

// AuthUsecase is a high-level interface for handling the generation
// and verification of access tokens.
type AuthUsecase interface {
	Token(ctx context.Context, cred *dto.UserCredential) result.Result
	Verify(ctx context.Context, tokenData []byte) result.Result
	VerifyWithScopes(ctx context.Context, tokenData []byte, scopes ...string) result.Result
}

type authUsecase struct {
	repo repository.UserRepository
	pwd password.Service
	auth *auth.Service
}

// NewAuthUsecase returns a new instance of AuthUsecase with the given repo.
func NewAuthUsecase(repo repository.UserRepository) AuthUsecase {
	return &authUsecase{
		repo: repo,
		pwd: password.New(),
		auth: auth.New(),
	}
}

// Token generates a new access token for the user with the
// given credentials. If the credentials are invalid a failed result
// will be returned, with a bad request status code.
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
	if t == nil {
		errMsg := "Oops, something went wrong when generating your access token :/"
		return result.Failure(errMsg)
	}

	ac := auth.NewAccessToken(t, exp)
	return result.Ok().WithValue(ac)
}

// Verify verifies the given access token.
func (u *authUsecase) Verify(ctx context.Context, tokenData []byte) result.Result {
	ok := u.auth.VerifyToken(ctx, tokenData)
	if !ok {
		return result.Failure("invalid token").
			WithStatusCode(http.StatusUnauthorized)
	}

	return result.Ok()
}

// VerifyWithScopes verifies the token, then verifies if the token contains
// any of the given scopes.
func (u *authUsecase) VerifyWithScopes(ctx context.Context, tokenData []byte, scopes ...string) result.Result {
	res := u.Verify(ctx, tokenData)
	if !res.IsOk() {
		return res
	}

	t := auth.Token(tokenData)
	tokenScopes := t.Strings(auth.ClaimTypeScopes)

	allowedScopes := make(map[string]int)
	for _, s := range scopes {
		allowedScopes[s] = 1
	}

	for _, ts := range tokenScopes {
		if _, ok := allowedScopes[ts]; ok {
			return result.Ok()
		}
	}

	return result.Failure("You're not allowed to access this resource :(").
		WithStatusCode(http.StatusForbidden)
}