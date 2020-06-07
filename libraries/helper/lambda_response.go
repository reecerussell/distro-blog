package helper

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/reecerussell/distro-blog/libraries/contextkey"

	"github.com/aws/aws-lambda-go/events"

	"github.com/reecerussell/distro-blog/libraries/result"
)

// Response creates a new APIGatewayProxyResponse form a given result.Result.
func Response(ctx context.Context, res result.Result, req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var resp events.APIGatewayProxyResponse
	var data responseWrapper
	success, status, value, err := res.Deconstruct()
	if !success {
		if status == 0 {
			resp.StatusCode = http.StatusInternalServerError
		} else {
			resp.StatusCode = status
		}

		msg := err.Error()
		data.ErrorMessage = &msg
	} else {
		if status == 0 {
			resp.StatusCode = http.StatusOK
		} else {
			resp.StatusCode = status
		}

		data.Data = value
	}

	jsonBytes, _ := json.Marshal(&data)
	resp.Body = string(jsonBytes)
	resp.IsBase64Encoded = false
	mapCORS(ctx, req, &resp)

	return resp
}

func mapCORS(ctx context.Context, req events.APIGatewayProxyRequest, resp *events.APIGatewayProxyResponse) {
	if resp.Headers == nil {
		resp.Headers = make(map[string]string)
	}

	if v := ctx.Value(contextkey.ContextKey("CORS_ORIGIN")); v != nil {
		resp.Headers["Access-Control-Allow-Origin"] = v.(string)
	} else {
		resp.Headers["Access-Control-Allow-Origin"] = "*"
	}

	if v := ctx.Value(contextkey.ContextKey("CORS_HEADERS")); v != nil {
		resp.Headers["Access-Control-Allow-Headers"] = v.(string)
	}

	resp.Headers["Access-Control-Allow-Method"] = req.HTTPMethod
}

type responseWrapper struct {
	ErrorMessage *string     `json:"error,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}
