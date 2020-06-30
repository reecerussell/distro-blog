package helper

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/reecerussell/distro-blog/auth"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/reecerussell/distro-blog/libraries/contextkey"
)

// PopulateContext populates the given context with the values from the API request type.
func PopulateContext(ctx context.Context, req events.APIGatewayProxyRequest) context.Context {
	for k, v := range req.StageVariables {
		ctx = context.WithValue(ctx, contextkey.ContextKey(k), v)
	}

	h ,ok := req.Headers["Authorization"]
	if !ok {
		h, ok = req.Headers["authorization"]
	}

	if ok {
		parts := strings.Split(h, ".")
		if len(parts) == 3 {
			bytes, _ := base64.StdEncoding.DecodeString(parts[1])
			var payload map[string]interface{}
			json.Unmarshal(bytes, &payload)

			if v, ok := payload[auth.ClaimTypeUserId]; ok {
				ctx = context.WithValue(ctx, contextkey.ContextKey("user_id"), v)
			}
		}
	}

	return ctx
}

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