package helper

import (
	"context"
	"encoding/base64"
	"encoding/json"

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

// TODO: test this.
func ReadBody(req events.APIGatewayProxyRequest, dst interface{}) {
	if req.IsBase64Encoded {
		data, _ := base64.StdEncoding.DecodeString(req.Body)
		json.Unmarshal(data, dst)
	} else {
		json.Unmarshal(string(req.Body), dst)
	}
}