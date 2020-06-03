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

func TestValidatePassword(t *testing.T) {
	testPasswords := []string{
		"Password_1234",
		"Mypassword#2",
	}

	for i, p := range testPasswords {
		err := serv.Validate(p)
		if err != nil {
			t.Errorf("validate[%d]: expected to be valid: %v", i, err)
		}
	}
}

func TestValidateEmptyPassword(t *testing.T) {
	err := serv.Validate("")
	if err == nil {
		t.Errorf("expected an error, but got none")
	}
}

func TestSetValidationOptions(t *testing.T) {
	opts := &Options{
		RequiredLength:         6,
		RequireUppercase:       true,
		RequireLowercase:       true,
		RequireNonAlphanumeric: false,
		RequireDigit:           true,
		RequiredUniqueChars:    0,
	}

	serv.SetValidationOptions(opts)
	s := serv.(*service)

	if s.validationOptions.RequireDigit != opts.RequireDigit {
		t.Errorf("RequireDigit: expected %v, but got %v", opts.RequireDigit, s.validationOptions.RequireDigit)
	}

	if s.validationOptions.RequiredUniqueChars != opts.RequiredUniqueChars {
		t.Errorf("RequiredUniqueChars: expected %d, but got %d", opts.RequiredUniqueChars, s.validationOptions.RequiredUniqueChars)
	}

	if s.validationOptions.RequiredLength != opts.RequiredLength {
		t.Errorf("RequiredLength: expected %d, but got %d", opts.RequiredLength, s.validationOptions.RequiredLength)
	}

	if s.validationOptions.RequireUppercase != opts.RequireUppercase {
		t.Errorf("RequireUppercase: expected %v, but got %v", opts.RequireUppercase, s.validationOptions.RequireUppercase)
	}

	if s.validationOptions.RequireLowercase != opts.RequireLowercase {
		t.Errorf("RequireLowercase: expected %v, but got %v", opts.RequireLowercase, s.validationOptions.RequireLowercase)
	}

	if s.validationOptions.RequireNonAlphanumeric != opts.RequireNonAlphanumeric {
		t.Errorf("RequireNonAlphanumeric: expected %v, but got %v", opts.RequireNonAlphanumeric, s.validationOptions.RequireNonAlphanumeric)
	}
}

func TestValidateRequiresUppercasePassword(t *testing.T) {
	opts := &Options{
		RequiredLength:         6,
		RequireUppercase:       true,
		RequireLowercase:       false,
		RequireNonAlphanumeric: false,
		RequireDigit:           false,
		RequiredUniqueChars:    0,
	}

	serv.SetValidationOptions(opts)

	pwd := "my test lowercase pass"
	err := serv.Validate(pwd)
	if err == nil {
		t.Errorf("expected an error but got nil")
	}

	opts.RequireUppercase = false
	err = serv.Validate(pwd)
	if err != nil {
		t.Errorf("expected nil but got: %v", err)
	}
}

func TestValidateRequiresLowercasePassword(t *testing.T) {
	opts := &Options{
		RequiredLength:         6,
		RequireUppercase:       false,
		RequireLowercase:       true,
		RequireNonAlphanumeric: false,
		RequireDigit:           false,
		RequiredUniqueChars:    0,
	}

	serv.SetValidationOptions(opts)

	pwd := "MY UPPERCASE PASSWORD"
	err := serv.Validate(pwd)
	if err == nil {
		t.Errorf("expected an error but got nil")
	}

	opts.RequireLowercase = false
	err = serv.Validate(pwd)
	if err != nil {
		t.Errorf("expected nil but got: %v", err)
	}
}

func TestValidateShortPassword(t *testing.T) {
	opts := &Options{
		RequiredLength:         6,
		RequireUppercase:       false,
		RequireLowercase:       true,
		RequireNonAlphanumeric: false,
		RequireDigit:           false,
		RequiredUniqueChars:    0,
	}

	serv.SetValidationOptions(opts)

	pwd := "short"
	err := serv.Validate(pwd)
	if err == nil {
		t.Errorf("expected an error but got nil")
	}

	pwd = "long password"
	err = serv.Validate(pwd)
	if err != nil {
		t.Errorf("expected nil but got: %v", err)
	}
}

func TestValidateRequiresNonAlphanumericPassword(t *testing.T) {
	opts := &Options{
		RequiredLength:         6,
		RequireUppercase:       false,
		RequireLowercase:       false,
		RequireNonAlphanumeric: true,
		RequireDigit:           false,
		RequiredUniqueChars:    0,
	}

	serv.SetValidationOptions(opts)

	pwd := "nonalpha"
	err := serv.Validate(pwd)
	if err == nil {
		t.Errorf("expected an error but got nil")
	}

	err = serv.Validate("non-alpha")
	if err != nil {
		t.Errorf("expected nil but got: %v", err)
	}

	opts.RequireNonAlphanumeric = false
	err = serv.Validate(pwd)
	if err != nil {
		t.Errorf("expected nil but got: %v", err)
	}
}

func TestValidateRequiresDigitPassword(t *testing.T) {
	opts := &Options{
		RequiredLength:         6,
		RequireUppercase:       false,
		RequireLowercase:       false,
		RequireNonAlphanumeric: false,
		RequireDigit:           true,
		RequiredUniqueChars:    0,
	}

	serv.SetValidationOptions(opts)

	pwd := "sdfghj"
	err := serv.Validate(pwd)
	if err == nil {
		t.Errorf("expected an error but got nil")
	}

	err = serv.Validate("3434hjk")
	if err != nil {
		t.Errorf("expected nil but got: %v", err)
	}

	opts.RequireDigit = false
	err = serv.Validate(pwd)
	if err != nil {
		t.Errorf("expected nil but got: %v", err)
	}
}

func TestValidateRequiresUniqueCharsPassword(t *testing.T) {
	opts := &Options{
		RequiredLength:         6,
		RequireUppercase:       false,
		RequireLowercase:       false,
		RequireNonAlphanumeric: false,
		RequireDigit:           false,
		RequiredUniqueChars:    2,
	}

	serv.SetValidationOptions(opts)

	pwd := "aaaaaa"
	err := serv.Validate(pwd)
	if err == nil {
		t.Errorf("expected an error but got nil")
	}

	err = serv.Validate("abcdef")
	if err != nil {
		t.Errorf("expected nil but got: %v", err)
	}

	opts.RequiredUniqueChars = 0
	err = serv.Validate(pwd)
	if err != nil {
		t.Errorf("expected nil but got: %v", err)
	}
}
