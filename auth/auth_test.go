package auth

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/reecerussell/distro-blog/libraries/contextkey"
)

var (
	testService *Service
	testKeyId = "alias/distro-jwt"
)

func init() {
	testService = New()
}

func TestAuthService_GenerateToken(t *testing.T) {
	claims := map[string]interface{}{
		"name": "John Doe",
	}
	userID := "139721"

	ctx := context.Background()
	ctx = context.WithValue(ctx, contextkey.ContextKey("JWT_KEY_ID"), testKeyId)

	tkn := testService.NewToken(ctx).
		SetExpiry(time.Now().UTC().Add(1 * time.Minute)).
		AddClaim(ClaimTypeUserId, userID).
		AddClaims(claims).
		Build()

	ok := testService.VerifyToken(ctx, tkn)
	if !ok {
		t.Errorf("expected token to be valid")
	}

	t.Run("No KeyID Set", func(t *testing.T) {
		ctx := context.Background()
		tkn := testService.NewToken(ctx).AddClaim(ClaimTypeUserId, userID).Build()
		if tkn != nil {
			t.Errorf("expected error and tkn to be nil")
		}
	})

	t.Run("Invalid Key", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, contextkey.ContextKey("JWT_KEY_ID"), "some random key id")
		tkn := testService.NewToken(ctx).AddClaim(ClaimTypeUserId, userID).Build()
		if tkn != nil {
			t.Errorf("expected error and tkn to be nil")
		}
	})
}

func TestTokenBuilder_AddClaims(t *testing.T) {
	claims := map[string]interface{}{
		"id": 5,
		"name": "John",
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, contextkey.ContextKey("JWT_KEY_ID"), testKeyId)

	tkn := testService.NewToken(ctx).AddClaims(claims).Build()
	payload := tkn.getPayload()

	if v := payload["id"].(float64); v != float64(claims["id"].(int)) {
		t.Errorf("expected 'id' claim to be: %v, but got: %v", claims["id"], v)
	}

	if v := payload["name"].(string); v != claims["name"].(string) {
		t.Errorf("expected 'name' claim to be: %v, but got: %v", claims["name"], v)
	}
}

func TestToken_Number(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextkey.ContextKey("JWT_KEY_ID"), testKeyId)

	tkn := testService.NewToken(ctx).AddClaim("id", 5).Build()

	v := tkn.Number("id")
	if v == nil {
		t.Errorf("unexpected a non-nil pointer")
	}

	if *v != float64(5) {
		t.Errorf("expected claim 'id' to be 5 but got: %v", *v)
	}

	t.Run("Non-existent Claim", func(t *testing.T) {
		v = tkn.Number("age")
		if v != nil {
			t.Errorf("expected a nil pointer")
		}
	})
}

func TestToken_Strings(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextkey.ContextKey("JWT_KEY_ID"), testKeyId)

	names := []string{"John", "Jane"}
	tkn := testService.NewToken(ctx).AddClaim("names", names).Build()

	v := tkn.Strings("names")
	if v == nil {
		t.Errorf("unexpected a non-nil slice")
	}

	for i, n := range v {
		if n != names[i] {
			t.Errorf("name[i] expected to be '%s' but was '%s'", names[i], n)
		}
	}

	t.Run("Non-existent Claim", func(t *testing.T) {
		v = tkn.Strings("age")
		if v != nil {
			t.Errorf("expected a nil pointer")
		}
	})
}

func TestToken_String(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextkey.ContextKey("JWT_KEY_ID"), testKeyId)

	tkn := testService.NewToken(ctx).AddClaim("id", 5).Build()

	if v := string(tkn); v != tkn.String() {
		t.Errorf("expected '%s' but got '%s'", v, tkn.String())
	}
}

func TestService_VerifyToken(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextkey.ContextKey("JWT_KEY_ID"), testKeyId)

	tkn := testService.NewToken(ctx).AddClaim("id", 5).Build()
	if !testService.VerifyToken(ctx, tkn) {
		t.Errorf("expected token to be valid")
	}

	t.Run("Malformed Header", func(t *testing.T) {
		mt := tkn[1:] // deform header
		if testService.VerifyToken(ctx, mt) {
			t.Errorf("expected token to be invalid")
		}
	})

	t.Run("Malformed Signature", func(t *testing.T) {
		mt := tkn[:len(tkn)-1]
		if testService.VerifyToken(ctx, mt) {
			t.Errorf("expected token to be invalid")
		}
	})

	t.Run("Malformed Structure", func(t *testing.T) {
		i := bytes.IndexByte(tkn, '.')
		if testService.VerifyToken(ctx, tkn[i+1:]) {
			t.Errorf("expected token to be invalid")
		}
	})
}

func TestTokenBuilder_SetExpiry(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextkey.ContextKey("JWT_KEY_ID"), testKeyId)

	exp := time.Now().UTC().Add(1 * time.Minute)
	tkn := testService.NewToken(ctx).SetExpiry(exp).Build()

	expected := convertToMilliseconds(exp)
	actual := tkn.Number(ClaimTypeExpiry)
	if actual == nil {
		t.Errorf("expected a non-nil pointer")
		return
	}

	if *actual != expected {
		t.Errorf("expected '%f' but got '%f'", expected, *actual)
	}
}

func TestTokenBuilder_SetNotBefore(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextkey.ContextKey("JWT_KEY_ID"), testKeyId)

	exp := time.Now().UTC()
	tkn := testService.NewToken(ctx).SetNotBefore(exp).Build()

	expected := convertToMilliseconds(exp)
	actual := tkn.Number(ClaimTypeNotBefore)
	if actual == nil {
		t.Errorf("expected a non-nil pointer")
		return
	}

	if *actual != expected {
		t.Errorf("expected '%f' but got '%f'", expected, *actual)
	}
}

func TestTokenBuilder_SetIssuedAt(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextkey.ContextKey("JWT_KEY_ID"), testKeyId)

	exp := time.Now().UTC()
	tkn := testService.NewToken(ctx).SetIssuedAt(exp).Build()

	expected := convertToMilliseconds(exp)
	actual := tkn.Number(ClaimTypeIssuedAt)
	if actual == nil {
		t.Errorf("expected a non-nil pointer")
		return
	}

	if *actual != expected {
		t.Errorf("expected '%f' but got '%f'", expected, *actual)
	}
}

// TODO: add test for ensuring payload is valid