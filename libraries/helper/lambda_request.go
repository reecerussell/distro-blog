package helper

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

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
func ReadBody(req events.APIGatewayProxyRequest, dst interface{}) error {
	defaultErr := fmt.Errorf("request contained invalid and/or malformed data")

	if req.IsBase64Encoded {
		data, err := base64.StdEncoding.DecodeString(req.Body)
		if err != nil {
			return defaultErr
		}

		err = json.Unmarshal(data, dst)
		if err != nil {
			return defaultErr
		}
	} else {
		err := json.Unmarshal([]byte(req.Body), dst)
		if err != nil {
			return defaultErr
		}
	}

	return nil
}