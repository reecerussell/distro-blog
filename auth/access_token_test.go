package auth

import (
	"context"
	"github.com/reecerussell/distro-blog/libraries/contextkey"
	"testing"
	"time"
)

func TestNewAccessToken(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextkey.ContextKey("JWT_KEY_ID"), "alias/distro-jwt")

	tkn := testService.NewToken(ctx).AddClaim("id", 5).Build()
	exp := time.Now().UTC().Add(time.Second * 30)

	ac := NewAccessToken(tkn, exp)
	if ac.Token != tkn.String() {
		t.Errorf("token mismatch")
	}

	if ac.Expires != exp.Unix() {
		t.Errorf("expiry mismatch")
	}
}
