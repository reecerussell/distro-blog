package helper

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/reecerussell/distro-blog/libraries/result"
)

func TestResponseSuccess(t *testing.T) {
	ctx := context.Background()
	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
	}
	res := result.Ok()

	resp := Response(ctx, res, req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected a status code of 200, but got: %d", resp.StatusCode)
	}

	// expected data
	var data responseWrapper
	jsonBytes, _ := json.Marshal(&data)
	json := string(jsonBytes)

	if resp.Body != json {
		t.Errorf("expected the response body to be '%s' but got '%s'", json, resp.Body)
	}

	// cors
	if v := resp.Headers["Access-Control-Allow-Origin"]; v != "*" {
		t.Errorf("expected '*' but got '%s'", v)
	}

	if _, ok := resp.Headers["Access-Control-Allow-Headers"]; ok {
		t.Errorf("expected to now have CORS header, header")
	}

	if v := resp.Headers["Access-Control-Allow-Method"]; v != req.HTTPMethod {
		t.Errorf("expected '%s' but got '%s'", req.HTTPMethod, v)
	}
}

func TestResponseSuccessWithStatusCode(t *testing.T) {
	ctx := context.Background()
	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
	}
	res := result.Ok().WithStatusCode(http.StatusNoContent)

	resp := Response(ctx, res, req)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected a status code of %d, but got: %d", http.StatusNoContent, resp.StatusCode)
	}
}

func TestResponseSuccessWithData(t *testing.T) {
	ctx := context.Background()
	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
	}
	testData := "Hello World"
	res := result.Ok().WithValue(testData)

	resp := Response(ctx, res, req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected a status code of 200, but got: %d", resp.StatusCode)
	}

	// expected data
	var data responseWrapper
	data.Data = testData
	jsonBytes, _ := json.Marshal(&data)
	json := string(jsonBytes)

	if resp.Body != json {
		t.Errorf("expected the response body to be '%s' but got '%s'", json, resp.Body)
	}
}

func TestResponseFailure(t *testing.T) {
	ctx := context.Background()
	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
	}
	testErrorMessage := "error"
	res := result.Failure(testErrorMessage)

	resp := Response(ctx, res, req)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected a status code of 500, but got: %d", resp.StatusCode)
	}

	// expected data
	var data responseWrapper
	data.ErrorMessage = &testErrorMessage
	jsonBytes, _ := json.Marshal(&data)
	json := string(jsonBytes)

	if resp.Body != json {
		t.Errorf("expected the response body to be '%s' but got '%s'", json, resp.Body)
	}

	// cors
	if v := resp.Headers["Access-Control-Allow-Origin"]; v != "*" {
		t.Errorf("expected '*' but got '%s'", v)
	}

	if _, ok := resp.Headers["Access-Control-Allow-Headers"]; ok {
		t.Errorf("expected to now have CORS header, header")
	}

	if v := resp.Headers["Access-Control-Allow-Method"]; v != req.HTTPMethod {
		t.Errorf("expected '%s' but got '%s'", req.HTTPMethod, v)
	}
}

func TestResponseFailureWithStatusCode(t *testing.T) {
	ctx := context.Background()
	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
	}
	res := result.Failure("").WithStatusCode(http.StatusBadRequest)

	resp := Response(ctx, res, req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a status code of %d, but got: %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestMapCORSWithKeys(t *testing.T) {
	req := events.APIGatewayProxyRequest{
		StageVariables: map[string]string{
			"CORS_ORIGIN":  "my web site",
			"CORS_HEADERS": "my header",
		},
		HTTPMethod: http.MethodGet,
	}
	ctx := PopulateContext(context.Background(), req)
	res := result.Ok()

	resp := Response(ctx, res, req)

	o, ok := resp.Headers["Access-Control-Allow-Origin"]
	if !ok {
		t.Errorf("expected CORS origin header")
	} else {
		if o != "my web site" {
			t.Errorf("expected 'my web site' but got '%s'", o)
		}
	}

	h, ok := resp.Headers["Access-Control-Allow-Headers"]
	if !ok {
		t.Errorf("expected CORS header, header")
	} else {
		if h != "my header" {
			t.Errorf("expected 'my header' but got '%s'", h)
		}
	}

	if v := resp.Headers["Access-Control-Allow-Method"]; v != req.HTTPMethod {
		t.Errorf("expected method '%s' but got '%s'", req.HTTPMethod, v)
	}
}

func TestMapCORSWithoutKeys(t *testing.T) {
	req := events.APIGatewayProxyRequest{
		StageVariables: map[string]string{},
		HTTPMethod:     http.MethodGet,
	}
	ctx := PopulateContext(context.Background(), req)
	res := result.Ok()

	resp := Response(ctx, res, req)

	o, ok := resp.Headers["Access-Control-Allow-Origin"]
	if !ok {
		t.Errorf("expected CORS origin header")
	} else {
		if o != "*" {
			t.Errorf("expected '*' but got '%s'", o)
		}
	}

	_, ok = resp.Headers["Access-Control-Allow-Headers"]
	if ok {
		t.Errorf("didn't expect CORS header, header")
	}

	if v := resp.Headers["Access-Control-Allow-Method"]; v != req.HTTPMethod {
		t.Errorf("expected method '%s' but got '%s'", req.HTTPMethod, v)
	}
}
