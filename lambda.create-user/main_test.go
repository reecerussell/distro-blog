package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/reecerussell/distro-blog/domain/dto"
)

func TestHandleCreateUser(t *testing.T) {
	data := &dto.CreateUser{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "handleCreateUser@test.com",
		Password:  "MyTestPassword123",
	}
	bytes, _ := json.Marshal(&data)
	encJSON := base64.StdEncoding.EncodeToString(bytes)

	ctx := context.Background()
	req := events.APIGatewayProxyRequest{
		HTTPMethod:      http.MethodPost,
		Body:            encJSON,
		IsBase64Encoded: true,
	}

	resp, _ := handleCreateUser(ctx, req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected: %d, Actual: %d", http.StatusOK, resp.StatusCode)
	}
}

func TestHandleCreateUserWithInvalidBody(t *testing.T) {
	ctx := context.Background()
	req := events.APIGatewayProxyRequest{
		HTTPMethod:      http.MethodPost,
		Body:            "invalid data",
		IsBase64Encoded: false,
	}

	resp, _ := handleCreateUser(ctx, req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a status code of '%d' but got '%d'", http.StatusBadRequest, resp.StatusCode)
	}
}
