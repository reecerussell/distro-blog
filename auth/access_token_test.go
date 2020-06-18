package auth

import (
	"testing"
	"time"
)

func TestNewAccessToken(t *testing.T) {
	tkn := testService.NewToken().AddClaim("id", 5).Build()
	exp := time.Now().UTC().Add(time.Second * 30)

	ac := NewAccessToken(tkn, exp)
	if ac.Token != tkn.String() {
		t.Errorf("token mismatch")
	}

	if ac.Expires != exp.Unix() {
		t.Errorf("expiry mismatch")
	}
}
