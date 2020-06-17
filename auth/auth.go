package auth

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var encoding = base64.RawURLEncoding

// Common errors.
var (
	ErrMalformedHeader = errors.New("malformed token header")
	ErrMalformedSignature = errors.New("malformed token signature")
	ErrMalformedStructure = errors.New("malformed token structure")
)

// Registered claims
const (
	ClaimTypeExpiry = "exp"
	ClaimTypeNotBefore = "nbf"
	ClaimTypeIssuedAt = "iat"
)

// Custom claim types.
const (
	ClaimTypeUserId = "uid"
	ClaimTypeScopes = "scp"
)

const (
	ScopeUserRead = "users:read"
	ScopeUserWrite = "users:write"
)

// Service is used to handle authentication and authorization of users
// using JSON-Web-Tokens.
type Service struct {
	sk       *rsa.PrivateKey
	alg crypto.Hash
}

func New(privKey *rsa.PrivateKey, alg crypto.Hash) *Service {
	return &Service{
		sk:       privKey,
		alg: alg,
	}
}

func (as *Service) VerifyToken(data []byte) bool {
	token := Token(data)
	ld, sig, err := token.scan()
	if err != nil {
		return false
	}

	digest := as.alg.New()
	digest.Write(data[:ld])

	err = rsa.VerifyPKCS1v15(&as.sk.PublicKey, as.alg, digest.Sum(sig[len(sig):]), sig)
	if err != nil {
		// signature mismatch
		return false
	}

	exp := token.Number(ClaimTypeExpiry)
	nbf := token.Number(ClaimTypeNotBefore)
	n := convertToMilliseconds(time.Now().UTC())

	// ensure token expiry is within bounds
	return (exp == nil || *exp > n) && (nbf == nil || *nbf <= n)
}

// TokenBuilder is an extension type to the Service struct, used to
// construct JSON-Web-Tokens.
type TokenBuilder struct {
	as *Service
	claims map[string]interface{}
}

// Token is a type used to provide extensions to a byte array. A token
// is a []byte is every way, just with extra methods.
type Token []byte

// NewToken returns a new builder, used to construct a token.
func (as *Service) NewToken() *TokenBuilder {
	return &TokenBuilder{
		as: as,
		claims: make(map[string]interface{}),
	}
}

// AddClaim adds a claim to the token, with the given name and value.
// If a claim with the given name already exists, the value will
// be overridden with the new value.
func (tb *TokenBuilder) AddClaim(name string, value interface{}) *TokenBuilder {
	tb.claims[name] = value

	return tb
}

// AddClaims adds a collection of claims to the token.
func (tb *TokenBuilder) AddClaims(claims map[string]interface{}) *TokenBuilder {
	for k, v := range claims {
		tb.AddClaim(k, v)
	}

	return tb
}

// SetExpiry sets the "Expiry" claim to the given time, in the form
// of the number of milliseconds since 1970-01-01T00:00:00Z UTC,
// ignoring leap seconds.
func (tb *TokenBuilder) SetExpiry(t time.Time) *TokenBuilder {
	ms := convertToMilliseconds(t)
	tb.AddClaim(ClaimTypeExpiry, ms)
	return tb
}

// SetIssuedAt sets the "Issued At" claim to the given time, in the form
// of the number of milliseconds since 1970-01-01T00:00:00Z UTC,
// ignoring leap seconds.
func (tb *TokenBuilder) SetIssuedAt(t time.Time) *TokenBuilder {
	ms := convertToMilliseconds(t)
	tb.AddClaim(ClaimTypeIssuedAt, ms)
	return tb
}

// SetNotBefore sets the "Not Before" claim to the given time, in the form
// of the number of milliseconds since 1970-01-01T00:00:00Z UTC,
// ignoring leap seconds.
func (tb *TokenBuilder) SetNotBefore(t time.Time) *TokenBuilder {
	ms := convertToMilliseconds(t)
	tb.AddClaim(ClaimTypeNotBefore, ms)
	return tb
}

func convertToMilliseconds(t time.Time) float64 {
	return float64(t.UnixNano() / 1e9)
}

// Build constructs the token using the data from the TokenBuilder.
func (tb *TokenBuilder) Build() Token {
	// encode the header
	headerPayload := map[string]string {
		"alg": "RSA256",
		"typ": "JWT",
	}
	headerData, _ := json.Marshal(headerPayload)
	header := make([]byte, encoding.EncodedLen(len(headerData)))
	encoding.Encode(header, headerData)

	// construct token
	payload, _ := json.Marshal(tb.claims)
	l := len(header) + 1 + encoding.EncodedLen(len(payload))
	token := make([]byte, l, l+1+encoding.EncodedLen(tb.as.sk.Size()))

	i := copy(token, header)
	token[i] = '.'
	i++
	encoding.Encode(token[i:], payload)
	token = token[:l]

	digest := tb.as.alg.New()
	digest.Write(token)

	// use signature space as a buffer while it's not set
	buf := token[len(token):]
	sig, _ := rsa.SignPKCS1v15(rand.Reader, tb.as.sk, tb.as.alg, digest.Sum(buf))

	i = len(token)
	token = token[:cap(token)]
	token[i] = '.'
	encoding.Encode(token[i+1:], sig)

	return token
}

// Number returns a number for the given claim from the Token payload.
// If there is not data for the claim, a nil pointer will be returned.
func (t *Token) Number(name string) *float64 {
	v, ok := t.getPayload()[name]
	if !ok {
		return nil
	}

	f := v.(float64)
	return &f
}

func (t *Token) getPayload() (payload map[string]interface{}) {
	encodedData := strings.Split(string(*t), ".")[1]
	rawData, _ := encoding.DecodeString(encodedData)
	_ = json.Unmarshal(rawData, &payload)
	return
}

func (t Token) scan() (int, []byte, error) {
	fd := bytes.IndexByte(t, '.')
	ld := bytes.LastIndexByte(t, '.')
	if ld <= fd {
		return 0, nil, ErrMalformedStructure
	}

	buf := make([]byte, encoding.DecodedLen(len(t)))
	_, err := encoding.Decode(buf, t[:fd])
	if err != nil {
		return 0, nil, ErrMalformedHeader
	}

	n, err := encoding.Decode(buf, t[ld+1:])
	if err != nil {
		return 0, nil, ErrMalformedSignature
	}

	return ld, buf[:n], nil
}