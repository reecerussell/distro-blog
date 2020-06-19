package usecase

import (
	"context"
	"os"
	"testing"

	"github.com/reecerussell/distro-blog/auth"
	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/libraries/contextkey"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/normalization"
	"github.com/reecerussell/distro-blog/password"
	"github.com/reecerussell/distro-blog/persistence"
)

var(
	testConnStringEmptySchema = os.Getenv("CONN_STRING_EMPTY_SCHEMA")
	testAuthUsecase AuthUsecase
)

func init() {
	db := database.NewMySQL(testConnString)
	repo := persistence.NewUserRepository(db)
	testAuthUsecase = NewAuthUsecase(repo)
}

func TestAuthUsecase_Token(t *testing.T) {
	email, password := "token@authUsecase.test", "MyTestPassword1"
	seedUser(email, password)

	ctx := context.Background()
	ctx = context.WithValue(ctx, contextkey.ContextKey("JWT_KEY_ID"), "alias/distro-jwt")

	d := &dto.UserCredential{
		Email: email,
		Password: password,
	}
	success, _, _, err := testAuthUsecase.Token(ctx, d).Deconstruct()
	if !success {
		t.Errorf("unexpected failure: %v", err)
	}

	t.Run("Unknown User", func(t *testing.T) {
		d := &dto.UserCredential{
			Email: "unknown email address",
			Password: password,
		}
		res := testAuthUsecase.Token(ctx, d)
		if res.IsOk() {
			t.Errorf("expecetd to fail")
		}
	})

	t.Run("Repository Failure", func(t *testing.T) {
		db := database.NewMySQL(testConnStringEmptySchema)
		repo := persistence.NewUserRepository(db)
		auth := NewAuthUsecase(repo)
		res := auth.Token(ctx, d)
		if res.IsOk() {
			t.Errorf("expected to fail")
		}
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		d := &dto.UserCredential{
			Email: email,
			Password: "invalid password",
		}
		res := testAuthUsecase.Token(ctx, d)
		if res.IsOk() {
			t.Errorf("expecetd to fail")
		}
	})

	t.Run("Token Build", func(t *testing.T) {
		res := testAuthUsecase.Token(context.Background(), d)
		if res.IsOk() {
			t.Errorf("expected to fail")
		}
	})
}

func TestAuthUsecase_Verify(t *testing.T) {
	email, password := "verify@authUsecase.test", "MyTestPassword1"
	seedUser(email, password)

	ctx := context.Background()
	ctx = context.WithValue(ctx, contextkey.ContextKey("JWT_KEY_ID"), "alias/distro-jwt")

	d := &dto.UserCredential{
		Email: email,
		Password: password,
	}
	success, _, value, err := testAuthUsecase.Token(ctx, d).Deconstruct()
	if !success {
		t.Errorf("unexpected failure: %v", err)
	}

	tokenData := []byte(value.(*auth.AccessToken).Token)

	// test start
	success, _, _, err = testAuthUsecase.Verify(ctx, tokenData).Deconstruct()
	if !success {
		t.Errorf("unexpected error: %v", err)
	}

	t.Run("No Key Set", func(t *testing.T) {
		res := testAuthUsecase.Verify(context.Background(), tokenData)
		if res.IsOk() {
			t.Errorf("expected to fail")
		}
	})

	t.Run("Deformed Token", func(t *testing.T) {
		data := tokenData[10:] //deform data
		res := testAuthUsecase.Verify(context.Background(), data)
		if res.IsOk() {
			t.Errorf("expected to fail")
		}
	})
}

func seedUser(email string, pwd string) {
	executeHelper("CALL create_user(UUID(), 'John', 'Doe', ?, ?, ?);",
		email, normalization.New().Normalize(email), password.New().Hash(pwd))
}