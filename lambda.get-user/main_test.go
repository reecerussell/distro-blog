package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

func TestHandleList(t *testing.T) {
	ctx := context.Background()
	req := events.APIGatewayProxyRequest{
		HTTPMethod:      http.MethodPost,
		PathParameters: map[string]string{
			"id": uuid.New().String(),
		},
	}

	// TODO: add more test cases
	resp, _ := handleGetUser(ctx, req)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected: %d, Actual: %d", http.StatusNotFound, resp.StatusCode)
	}
}