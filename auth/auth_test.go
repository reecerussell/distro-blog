package auth

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"
)

var (
	testPrivateKeyData = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQBmo5UHZEz3lHhuFIDQgDN3Yr5ic3A+We+UGKWaC/uCJsoD+lBH
WqpYv+OrrWoAH5Y3ROTra447MUOwuyPbeJECNhjA1KD8aEPi+bTMI4lreKaqMnCO
gjwpUXvbgg81ERz1DzYT5gXB5nFzf4R2LVJmUgDCN5srQE+xcvkmejE8CpaCF+Zq
h0LUAdJiQUgaBeAkmwrzcnqs2GBfV68iRJL++nfA9AQ/kH2Q17Ga3ipwLo1JbCTk
Zh4IfLdOpyMvVHwtqHR3lgGnCzFcmX0hl6AJ9oZkEFph2yPyrW2HYJPekBXIQ2Gk
7ziR8sZdwQ/fafS2v4fZS803VutWff/TPozjAgMBAAECggEAYNe/6bWNmZyQ9OyL
ji8oYGDe2e2p3mrlTori1bKwoGERAyfPT0QQrqR/oKCC/5LOHV/3ztkw3lDhWYN5
lb7ws3FvcaIuM3n9c8+/800kgC7asoPdB9mCAkpL3xWcW6nF9MNhduz2SbmxGhUb
WpXwxXJiHN5yniCUEQ42X0Oz6L2ru628sFioGWNumGm0bgkQtC69oKwpDNNYPxai
OrwIpXez8EymVunajaE2p68kW3vBxwEK9RU6OrsTSeozB4IgksKn05wuoxi67keX
nRdOVryamp2YpR45oxhM5YPuDoGR7kPk3o7VWrmO6AG4DEngSUX8TR/8tpRtQz/y
/TWVgQKBgQDJB3hOCq+9xRkyFY7DxuKUIH2O/VgNfbM6NiaQN87mr4OTK5CkOXeg
8PFd0IukSJdt0BbinJ6zgPgXD7MTwtzccbHAUIv+cWxrmYLnZu6UH7O+nIEIBc8r
wjh3ClaOToch8EpOr26/RdnQInFgG9q96h2NOGOo0V4HpAevlL2dswKBgQCCtI9/
V/ylMU487JDQXnmphlyA+z73lobylUy8vi+C3wmf6NTMuZsh+tCQ/ftecxpnpbO1
aa0nqRihWfCco8XNlcNbcrfMuvJkS12zsyWHqLRXxXjzyoXSl7/uAJI+2ZRE+LAU
2m12n+ltltXfGbZ5rAAQ58P79wQVqObtpDqcEQKBgDk1EP1UeTKN1l+0Vs5L2MrC
fDimy9n6/XgBVPQRjaWEKPNGoIC7gdmg9271G+gCaGVtpDWU0GzQtMkLRLDI8UUi
ba0GvvAHowzzwJbNafNpGiOSMf3weUZAnQTzQjJ5EmeME/lUXzW7UQKz6oOpKZSF
/Sbk9ydhfVq7SRykPVmVAoGAPYMygXDsQuY4du2ynY3I3iKQyFb15FmgOuxOyAkN
nR7QjcRq2cqEGvLKU2JkeafcBmlycO9CAYdQQydr2JwuzDkuToxnud9FkjPx7k9i
WzznWuNhsAJhBqJKPn1gVlnZsLgFTlsZ5xkNJ3k0QCH+wbZT9aDNmHhBINxzieWf
e7ECgYEAvxRcMlL0PLPUfCDRKmrDU6gs4QfiPNHWYlvqTGnTjyPqVAQScRBW2eCB
eTTr/63iG5wzMQL3PLGF+8KXURC2tuAlMeEPMMvFyc2aNB8FqHfx+W8/71YsIDR5
RlBmbI6m4YYLpjXGKO+XE5rgm3qLHeSK4r4MA31wWskw0QOG7ts=
-----END RSA PRIVATE KEY-----`
	testService *Service
)

func init() {
	block, _ := pem.Decode([]byte(testPrivateKeyData))
	pk, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	testService = New(pk, crypto.SHA512)
}

func TestAuthService_GenerateToken(t *testing.T) {
	claims := map[string]interface{}{
		"name": "John Doe",
	}
	userID := "139721"

	tkn := testService.NewToken().
		SetExpiry(time.Now().UTC().Add(1 * time.Minute)).
		AddClaim(ClaimTypeUserId, userID).
		AddClaims(claims).
		Build()

	ok := testService.VerifyToken(tkn)
	if !ok {
		t.Errorf("expected token to be valid")
	}
}

func TestTokenBuilder_AddClaims(t *testing.T) {
	claims := map[string]interface{}{
		"id": 5,
		"name": "John",
	}
	tkn := testService.NewToken().AddClaims(claims).Build()
	payload := tkn.getPayload()

	if v := payload["id"].(float64); v != float64(claims["id"].(int)) {
		t.Errorf("expected 'id' claim to be: %v, but got: %v", claims["id"], v)
	}

	if v := payload["name"].(string); v != claims["name"].(string) {
		t.Errorf("expected 'name' claim to be: %v, but got: %v", claims["name"], v)
	}
}

func TestToken_Number(t *testing.T) {
	tkn := testService.NewToken().AddClaim("id", 5).Build()

	v := tkn.Number("id")
	if v == nil {
		t.Errorf("expected a non-nil pointer")
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

func TestService_VerifyToken(t *testing.T) {
	tkn := testService.NewToken().AddClaim("id", 5).Build()
	if !testService.VerifyToken(tkn) {
		t.Errorf("expected token to be valid")
	}

	t.Run("Malformed Header", func(t *testing.T) {
		mt := tkn[1:] // deform header
		if testService.VerifyToken(mt) {
			t.Errorf("expected token to be invalid")
		}
	})

	t.Run("Malformed Signature", func(t *testing.T) {
		mt := tkn[:len(tkn)-1]
		if testService.VerifyToken(mt) {
			t.Errorf("expected token to be invalid")
		}
	})

	t.Run("Malformed Structure", func(t *testing.T) {
		i := bytes.IndexByte(tkn, '.')
		if testService.VerifyToken(tkn[i+1:]) {
			t.Errorf("expected token to be invalid")
		}
	})

	t.Run("Signature Mismatch", func(t *testing.T) {
		block, _ := pem.Decode([]byte(testPrivateKeyData))
		pk, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
		s := New(pk, crypto.SHA256)

		if s.VerifyToken(tkn) {
			t.Errorf("expected token to be invalid")
		}
	})
}

func TestTokenBuilder_SetExpiry(t *testing.T) {
	exp := time.Now().UTC().Add(1 * time.Minute)
	tkn := testService.NewToken().SetExpiry(exp).Build()

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
	exp := time.Now().UTC()
	tkn := testService.NewToken().SetNotBefore(exp).Build()

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
	exp := time.Now().UTC()
	tkn := testService.NewToken().SetIssuedAt(exp).Build()

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