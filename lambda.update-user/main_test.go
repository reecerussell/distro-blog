package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"

	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/libraries/normalization"
)

var testConnString string

func init() {
	testConnString = os.Getenv("CONN_STRING")
}

func TestHandler(t *testing.T) {
	id := seedUser("updateUserHandler@lambda.test")

	data := &dto.UpdateUser{
		ID: id,
		Firstname: "Jane",
		Lastname:  "Doe",
		Email:     "updateUserHandler@lambda.test",
	}
	bytes, _ := json.Marshal(&data)
	encJSON := base64.StdEncoding.EncodeToString(bytes)

	ctx := context.Background()
	req := events.APIGatewayProxyRequest{
		HTTPMethod:      http.MethodPut,
		Body:            encJSON,
		IsBase64Encoded: true,
	}

	resp, _ := handleUpdateUser(ctx, req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected: %d, Actual: %d", http.StatusOK, resp.StatusCode)
	}

	t.Run("Invalid Data", func(t *testing.T) {
		req := events.APIGatewayProxyRequest{
			HTTPMethod:      http.MethodPost,
			Body:            "invalid data",
			IsBase64Encoded: false,
		}

		resp, _ := handleUpdateUser(ctx, req)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected a status code of '%d' but got '%d'", http.StatusBadRequest, resp.StatusCode)
		}
	})
}

func seedUser(email string) string {
	query := "CALL `create_user`(?,?,?,?,?,?);"
	args := []interface{}{
		uuid.New().String(),
		"John",
		"Doe",
		email,
		normalization.New().Normalize(email),
		"egwuekhr",
	}

	db, err := sql.Open("mysql", testConnString)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		panic(err)
	}

	return args[0].(string)
}