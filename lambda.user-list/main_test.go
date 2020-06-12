package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandleList(t *testing.T) {
	ctx := context.Background()
	req := events.APIGatewayProxyRequest{
		HTTPMethod:      http.MethodPost,
	}

	resp, _ := handleList(ctx, req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected: %d, Actual: %d", http.StatusOK, resp.StatusCode)
	}
}