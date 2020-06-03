package password

import (
	"testing"
)

var serv = New()
var testPassword = "password_1234"

func TestHashPassword(t *testing.T) {
	hash := serv.Hash(testPassword)

	valid := serv.Verify(testPassword, hash)
	if !valid {
		t.Errorf("expected to be valid, but wasn't")
		return
	}
}

func TestVerifyWithInvalidHash(t *testing.T) {
	invalidHashes := []string{
		"",
		"rtyuykj",
		"3479032oed  f84",
	}

	for i, s := range invalidHashes {
		valid := serv.Verify(testPassword, s)
		if valid {
			t.Errorf("verify[%d]: expected false, but got true", i)
		}
	}
}

func TestInvalidAlgKey(t *testing.T) {
	hash := serv.Hash(testPassword)
	bytes, _ := encoding.DecodeString(hash)

	// change salt size, in header
	writeHeader(bytes, 1, 348) // ensure 348 is not a valid alg key.
	hash = encoding.EncodeToString(bytes)

	valid := serv.Verify(testPassword, hash)
	if valid {
		t.Errorf("expected password to be invalid")
	}
}

func TestInvalidSaltSize(t *testing.T) {
	hash := serv.Hash(testPassword)
	bytes, _ := encoding.DecodeString(hash)

	// change salt size, in header
	writeHeader(bytes, 9, 1)
	hash = encoding.EncodeToString(bytes)

	valid := serv.Verify(testPassword, hash)
	if valid {
		t.Errorf("expected password to be invalid")
	}
}

func TestInvalidKeySize(t *testing.T) {
	hash := serv.Hash(testPassword)
	bytes, _ := encoding.DecodeString(hash)

	// deform hash
	hash = encoding.EncodeToString(bytes[:len(bytes)-10])

	valid := serv.Verify(testPassword, hash)
	if valid {
		t.Errorf("expected password to be invalid")
	}
}
