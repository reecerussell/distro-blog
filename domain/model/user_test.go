package model

import (
	"testing"

	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/libraries/normalization"
	"github.com/reecerussell/distro-blog/password"
)

var testPasswordService = password.New()
var testNormalizer = normalization.New()
var longString = `PGcDxEDyVoe4KTTCDlC8GKbiPtL6MCxfQIlUiOi03AhiObscUSQU1dsZbOUp3VMXpGpLO7bDyFcM1H2XS3J1WlUTsi51SIONukqdb`

func TestUpdateFirstname(t *testing.T) {
	u := new(User)

	// Empty firstname - should fail.
	err := u.UpdateFirstname("")
	if err == nil {
		t.Errorf("empty: expected an error, but got none")
	}

	// Out of bounds firstname - should fail.
	err = u.UpdateFirstname(longString)
	if err == nil {
		t.Errorf("out of bounds: expected an error, but got none")
	}

	// valid name - should work.
	err = u.UpdateFirstname("John")
	if err != nil {
		t.Errorf("valid firstname: expected no error, but got: %v", err)
	}

	if u.firstname != "John" {
		t.Errorf("expected firstname to be 'John' but got: %s", u.firstname)
	}

	// no change - should work.
	u.firstname = "John" // seed firstname.
	err = u.UpdateFirstname(u.firstname)
	if err != nil {
		t.Errorf("no change: expected no error, but got: %v", err)
	}

	if u.firstname != "John" {
		t.Errorf("expected firstname to be 'John' but got: %s", u.firstname)
	}
}

func TestUpdateLastname(t *testing.T) {
	u := new(User)

	// Empty lastname - should fail.
	err := u.UpdateLastname("")
	if err == nil {
		t.Errorf("empty: expected an error, but got none")
	}

	// Out of bounds lastname - should fail.
	err = u.UpdateLastname(longString)
	if err == nil {
		t.Errorf("out of bounds: expected an error, but got none")
	}

	// valid name - should work.
	err = u.UpdateLastname("Doe")
	if err != nil {
		t.Errorf("valid lastname: expected no error, but got: %v", err)
	}

	if u.lastname != "Doe" {
		t.Errorf("expected lastname to be 'Doe' but got: %s", u.lastname)
	}

	// no change - should work.
	u.lastname = "Doe" // seed lastname.
	err = u.UpdateLastname(u.lastname)
	if err != nil {
		t.Errorf("no change: expected no error, but got: %v", err)
	}

	if u.lastname != "Doe" {
		t.Errorf("expected lastname to be 'Doe' but got: %s", u.lastname)
	}
}

func TestUpdateEmail(t *testing.T) {
	u := new(User)

	// empty email - should fail
	err := u.UpdateEmail("", testNormalizer)
	if err == nil {
		t.Errorf("empty: expected an error but got none")
	}

	// out of bounds - should fail
	err = u.UpdateEmail(longString, testNormalizer)
	if err == nil {
		t.Errorf("out of bounds: expected an error but got none")
	}

	// invalid email addresses
	invalidAddresses := []string{
		"my email address",
		"john@doe",
		"johndoe.com",
	}

	for i, a := range invalidAddresses {
		err = u.UpdateEmail(a, testNormalizer)
		if err == nil {
			t.Errorf("invalid[%d]: expected error but got none", i)
		}
	}

	// valid email
	err = u.UpdateEmail("john@doe.com", testNormalizer)
	if err != nil {
		t.Errorf("valid email: expected no error, but got: %v", err)
	}

	if u.email != "john@doe.com" {
		t.Errorf("expected email to be 'john@doe.com' but got: '%s'", u.email)
	}

	// no change
	u.email = "john@doe.com" // seed email.
	err = u.UpdateEmail(u.email, testNormalizer)
	if err != nil {
		t.Errorf("no change: expected no error, but got: %v", err)
	}

	if u.email != "john@doe.com" {
		t.Errorf("expected email to be 'john@doe.com' but got: '%s'", u.email)
	}
}

func TestSetPassword(t *testing.T) {
	pwdOpts := &password.Options{
		RequiredLength:         6,
		RequireUppercase:       true,
		RequireLowercase:       true,
		RequireNonAlphanumeric: false,
		RequireDigit:           true,
		RequiredUniqueChars:    0,
	}

	testPasswords := map[string]bool{
		"my password":    false,
		"Password1234":   true,
		"1234":           false,
		"My S3cur3 P1ss": true,
	}

	testPasswordService.SetValidationOptions(pwdOpts)

	u := new(User)

	for pwd, res := range testPasswords {
		err := u.setPassword(pwd, testPasswordService)
		if err == nil && !res {
			t.Errorf("expected an error but got none")
		} else if err != nil && res {
			t.Errorf("expected no error, but got: %v", err)
		}
	}
}

func TestNewUser(t *testing.T) {
	cu := &dto.CreateUser{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john@doe.com",
		Password:  "MyPassword123",
	}

	u, err := NewUser(cu, testPasswordService, testNormalizer)
	if err != nil {
		t.Errorf("valid user: expected no error, but got: %v", err)
	}

	if u == nil {
		t.Errorf("expected an instance of a user but got nil")
	}
}

func TestNewUserWithInvalidFirstname(t *testing.T) {
	cu := &dto.CreateUser{
		Firstname: "",
		Lastname:  "Doe",
		Email:     "john@doe.com",
		Password:  "MyPassword123",
	}

	_, err := NewUser(cu, testPasswordService, testNormalizer)
	if err == nil {
		t.Errorf("expected no error, but got: nil")
	}
}

func TestNewUserWithInvalidLastname(t *testing.T) {
	cu := &dto.CreateUser{
		Firstname: "John",
		Lastname:  "",
		Email:     "john@doe.com",
		Password:  "MyPassword123",
	}

	_, err := NewUser(cu, testPasswordService, testNormalizer)
	if err == nil {
		t.Errorf("expected no error, but got: nil")
	}
}

func TestNewUserWithInvalidEmail(t *testing.T) {
	cu := &dto.CreateUser{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "johndoe.com",
		Password:  "MyPassword123",
	}

	_, err := NewUser(cu, testPasswordService, testNormalizer)
	if err == nil {
		t.Errorf("expected no error, but got: nil")
	}
}

func TestNewUserWithInvalidPassword(t *testing.T) {
	cu := &dto.CreateUser{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john@doe.com",
		Password:  "",
	}

	_, err := NewUser(cu, testPasswordService, testNormalizer)
	if err == nil {
		t.Errorf("expected no error, but got: nil")
	}
}

func TestGetEmail(t *testing.T) {
	testEmail := "john@doe.com"
	u := &User{
		email: testEmail,
	}

	if v := u.Email(); v != testEmail {
		t.Errorf("expected '%s' but got '%s'", testEmail, v)
	}
}

func TestGetNormalizedEmail(t *testing.T) {
	testEmail := "JOHN@DOE.COM"
	u := &User{
		normalizedEmail: testEmail,
	}

	if v := u.NormalizedEmail(); v != testEmail {
		t.Errorf("expected '%s' but got '%s'", testEmail, v)
	}
}

func TestDataModel(t *testing.T) {
	u := &User{
		id:              "id",
		firstname:       "firstname",
		lastname:        "lastname",
		email:           "email",
		normalizedEmail: "normalized email",
		passwordHash:    "password hash",
	}

	dm := u.DataModel()

	if dm.ID != u.id {
		t.Errorf("expected id to be '%s' but got '%s'", u.id, dm.ID)
	}

	if dm.Firstname != u.firstname {
		t.Errorf("expected firstname to be '%s' but got '%s'", u.firstname, dm.Firstname)
	}

	if dm.Lastname != u.lastname {
		t.Errorf("expected lastname to be '%s' but got '%s'", u.lastname, dm.Lastname)
	}

	if dm.Email != u.email {
		t.Errorf("expected email to be '%s' but got '%s'", u.email, dm.Email)
	}

	if dm.NormalizedEmail != u.normalizedEmail {
		t.Errorf("expected normalized email to be '%s' but got '%s'", u.normalizedEmail, dm.NormalizedEmail)
	}

	if dm.PasswordHash != u.passwordHash {
		t.Errorf("expected password hash to be '%s' but got '%s'", u.passwordHash, dm.PasswordHash)
	}
}
