package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	authMod "github.com/reecerussell/distro-blog/auth"
	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/libraries/contextkey"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/normalization"
	"github.com/reecerussell/distro-blog/password"
	"github.com/reecerussell/distro-blog/persistence"
	"github.com/reecerussell/distro-blog/usecase"
	"os"
	"testing"
)

var (
	testConnString = os.Getenv("CONN_STRING")
	testAuth usecase.AuthUsecase
)

func init() {
	db := database.NewMySQL(testConnString)
	repo := persistence.NewUserRepository(db)
	testAuth = usecase.NewAuthUsecase(repo)
}

func TestHandleAuthentication(t *testing.T) {
	email, password, scope := "handleAuthentication@authorizer.lambda", "MyTestPassword", "users:read"
	methodArn := "arn:aws:execute-api:<region>:<account id>:<rest api id>/*/GET/users"
	seedUser(email, password, scope)

	ctx := context.Background()
	ctx = context.WithValue(ctx, contextkey.ContextKey("JWT_KEY_ID"), "alias/distro-jwt")
	d := &dto.UserCredential{
		Email: email,
		Password: password,
	}
	success, _, value, err := testAuth.Token(ctx, d).Deconstruct()
	if !success {
		t.Errorf("unexpected failure: %v", err)
	}

	tokenData := []byte(value.(*authMod.AccessToken).Token)
	token := "Bearer " + string(tokenData)
	authReq := events.APIGatewayCustomAuthorizerRequest{
		AuthorizationToken: token,
		MethodArn: methodArn,
	}

	_, err = handleAuthorization(context.Background(), authReq)
	if err != nil {
		t.Errorf("unpexpected error: %v", err)
		return
	}

	t.Run("Invalid Scheme", func(t *testing.T) {
		req := authReq
		req.AuthorizationToken = "Basic" + token[7:]
		_, err = handleAuthorization(context.Background(), req)
		if err == nil {
			t.Errorf("should've failed")
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := authReq
		req.AuthorizationToken = token[:len(token)-10] // deformed
		_, err = handleAuthorization(context.Background(), req)
		if err == nil {
			t.Errorf("should've failed")
		}
	})
}

func seedUser(email, pwd, scope string) {
	userID := uuid.New().String()
	executeHelper("CALL create_user(?, 'John', 'Doe', ?, ?, ?);",
		userID, email, normalization.New().Normalize(email), password.New().Hash(pwd))
	scopeID := uuid.New().String()
	executeHelper("INSERT INTO `scopes` (`id`,`name`,`description`) VALUES (?, ?, 'Test Scope')", scopeID, scope)
	executeHelper("INSERT INTO `user_scopes` (`user_id`,`scope_id`) VALUES (?,?);", userID, scopeID)
}

func executeHelper(query string, args ...interface{}) {
	db, err := sql.Open("mysql", testConnString)
	if err != nil {
		panic(fmt.Errorf("open: %v", err))
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		panic(fmt.Errorf("exec: %v", err))
	}
}
