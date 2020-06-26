package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/libraries/normalization"
	"github.com/reecerussell/distro-blog/password"
)

func TestHandleToken(t *testing.T) {
	// seed user
	email, pwd := "handleToken@lambda.test", "MySecurePass123"
	executeHelper("CALL create_user(UUID(), 'John','Doe', ?,?,?)",
		email, normalization.New().Normalize(email), password.New().Hash(pwd))

	cred := &dto.UserCredential{
		Email: email,
		Password: pwd,
	}
	bytes, _ := json.Marshal(cred)
	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		StageVariables: map[string]string{
			"JWT_KEY_ID": os.Getenv("JWT_KEY_ID"),
		},
		Body: string(bytes),
	}

	ctx := context.Background()
	resp, _ := handleToken(ctx, req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected a status code of 200 but got: %d", resp.StatusCode)
	}

	t.Run("Invalid Request Body", func(t *testing.T) {
		r := req
		r.Body = ""
		resp , _ := handleToken(ctx, r)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected a status code 400, but got: %d", resp.StatusCode)
		}
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		c := cred
		c.Password = "My Invalid Password"
		bytes, _ := json.Marshal(c)
		r := req
		r.Body = string(bytes)

		resp , _ := handleToken(ctx, r)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected a status code 400, but got: %d", resp.StatusCode)
		}
	})

	t.Run("Missing JWT Key ID", func(t *testing.T) {
		r := req
		r.StageVariables = nil

		resp , _ := handleToken(ctx, r)
		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected a status code 500, but got: %d", resp.StatusCode)
		}
	})
}

func executeHelper(query string, args ...interface{}) {
	db, err := sql.Open("mysql", os.Getenv("CONN_STRING"))
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		panic(err)
	}
}