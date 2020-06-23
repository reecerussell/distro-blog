package helper

import (
	"context"
	"encoding/base64"
	"encoding/json"
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

func TestReadBody(t *testing.T) {
	encoding := base64.StdEncoding
	testPayload := map[string]string{
		"Hello": "World",
	}

	t.Run("Valid Base64", func(t *testing.T) {
		bytes, _ := json.Marshal(testPayload)
		body := make([]byte, encoding.EncodedLen(len(bytes)))
		encoding.Encode(body, bytes)

		req := events.APIGatewayProxyRequest{
			Body: string(body),
			IsBase64Encoded: true,
		}

		var out map[string]string
		err := ReadBody(req, &out)
		if err != nil {
			t.Errorf("expected to be valid: %v", err)
		}
	})

	t.Run("Valid JSON", func(t *testing.T) {
		bytes, _ := json.Marshal(testPayload)
		req := events.APIGatewayProxyRequest{
			Body: string(bytes),
			IsBase64Encoded: false,
		}

		var out map[string]string
		err := ReadBody(req, &out)
		if err != nil {
			t.Errorf("expected to be valid: %v", err)
		}
	})

	t.Run("Invalid Base64", func(t *testing.T) {
		bytes, _ := json.Marshal(testPayload)
		body := make([]byte, encoding.EncodedLen(len(bytes)))
		encoding.Encode(body, bytes)

		req := events.APIGatewayProxyRequest{
			Body: string(body[10:]), // deformed base64
			IsBase64Encoded: true,
		}

		var out map[string]string
		err := ReadBody(req, &out)
		if err == nil {
			t.Errorf("expected to fail")
		}
	})

	t.Run("Valid Base64, Invalid JSON", func(t *testing.T) {
		bytes, _ := json.Marshal(testPayload)
		body := make([]byte, encoding.EncodedLen(len(bytes[3:])))
		encoding.Encode(body, bytes[3:])

		req := events.APIGatewayProxyRequest{
			Body: string(body),
			IsBase64Encoded: true,
		}

		var out map[string]string
		err := ReadBody(req, &out)
		if err == nil {
			t.Errorf("expected to fail")
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		bytes, _ := json.Marshal(testPayload)
		req := events.APIGatewayProxyRequest{
			Body: string(bytes[3:]), // deform json
			IsBase64Encoded: false,
		}

		var out map[string]string
		err := ReadBody(req, &out)
		if err == nil {
			t.Errorf("expected to fail")
		}
	})
}