package model

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"

	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/normalization"
	"github.com/reecerussell/distro-blog/password"
)

// User is a domain model for user records.
type User struct {
	id              string
	firstname       string
	lastname        string
	email           string
	normalizedEmail string
	passwordHash    string
}

// NewUser returns a new instance of a User domain model, after going
// through model validation.
func NewUser(data *dto.CreateUser, serv password.Service, norm normalization.Normalizer) (*User, error) {
	u := &User{
		id: uuid.New().String(),
	}

	err := u.UpdateFirstname(data.Firstname)
	if err != nil {
		return nil, err
	}

	err = u.UpdateLastname(data.Lastname)
	if err != nil {
		return nil, err
	}

	err = u.UpdateEmail(data.Email, norm)
	if err != nil {
		return nil, err
	}

	err = u.setPassword(data.Password, serv)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// UpdateFirstname updates the User's firstname.
func (u *User) UpdateFirstname(firstname string) error {
	l := len(firstname)

	switch true {
	case l < 1:
		return fmt.Errorf("firstname is required")
	case u.firstname == firstname:
		return nil
	case l > 45:
		return fmt.Errorf("firstname cannot be greater than 45 characters long")
	}

	u.firstname = firstname

	return nil
}

// UpdateLastname updates the User's lastname.
func (u *User) UpdateLastname(lastname string) error {
	l := len(lastname)

	switch true {
	case l < 1:
		return fmt.Errorf("lastname is required")
	case u.lastname == lastname:
		return nil
	case l > 45:
		return fmt.Errorf("lastname cannot be greater than 45 characters long")
	}

	u.lastname = lastname

	return nil
}

// UpdateEmail updates the User's email.
func (u *User) UpdateEmail(email string, normalizer normalization.Normalizer) error {
	l := len(email)
	re := regexp.MustCompile("[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,6}")

	switch true {
	case l < 1:
		return fmt.Errorf("email is required")
	case u.normalizedEmail == normalizer.Normalize(email):
		return nil
	case l > 100:
		return fmt.Errorf("email cannot be greater than 100 characters")
	case !re.MatchString(email):
		return fmt.Errorf("email is invalid")
	}

	u.email = email
	u.normalizedEmail = normalizer.Normalize(email)

	return nil
}

// Sets the user's password after validating and hashing it.
func (u *User) setPassword(password string, serv password.Service) error {
	err := serv.Validate(password)
	if err != nil {
		return err
	}

	u.passwordHash = serv.Hash(password)

	return nil
}