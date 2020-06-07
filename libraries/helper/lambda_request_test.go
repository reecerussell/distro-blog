package helper

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/reecerussell/distro-blog/libraries/contextkey"
)

func TestPopulateContext(t *testing.T) {
	req := events.APIGatewayProxyRequest{
		StageVariables: map[string]string{
			"Hello": "World",
		},
	}
	ctx := context.Background()
	ctx = PopulateContext(ctx, req)

	if v := ctx.Value(contextkey.ContextKey("Hello")); v != "World" {
		t.Errorf("expected 'World' but got '%v'", v)
	}
}
