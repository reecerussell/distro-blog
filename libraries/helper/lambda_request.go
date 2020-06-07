package helper

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/reecerussell/distro-blog/libraries/contextkey"
)

// PopulateContext populates the given context with the values from the API request type.
func PopulateContext(ctx context.Context, req events.APIGatewayProxyRequest) context.Context {
	for k, v := range req.StageVariables {
		ctx = context.WithValue(ctx, contextkey.ContextKey(k), v)
	}

	return ctx
}
