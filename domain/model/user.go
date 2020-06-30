package model

import (
	"context"
	"fmt"
	"github.com/reecerussell/distro-blog/domain/event"
	"github.com/reecerussell/distro-blog/domain/handler"
	"github.com/reecerussell/distro-blog/libraries/contextkey"
	"github.com/reecerussell/distro-blog/libraries/domainevents"
	"github.com/reecerussell/distro-blog/libraries/result"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"

	"github.com/reecerussell/distro-blog/domain/datamodel"
	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/libraries/normalization"
	"github.com/reecerussell/distro-blog/password"
)

// Common Audit Messages
const (
	AuditUserCreated = "USER_CREATED"
	AuditUserUpdated = "USER_UPDATED"
	AuditUserPasswordReset = "USER_PASSWORD_RESET"
	AuditUserPasswordChanged = "USER_PASSWORD_CHANGED"
)

func init() {
	domainevents.RegisterEventHandler(&event.AddUserAudit{}, &handler.AddUserAudit{})
}

// User is a domain model for user records.
type User struct {
	domainevents.Aggregate

	id              string
	firstname       string
	lastname        string
	email           string
	normalizedEmail string
	passwordHash    string

	scopes []*Scope
}

// NewUser returns a new instance of a User domain model, after going
// through model validation.
func NewUser(ctx context.Context, data *dto.CreateUser, serv password.Service, norm normalization.Normalizer) (*User, error) {
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

	var performingUserID string
	uid := ctx.Value(contextkey.ContextKey("user_id"))
	if uid != nil {
		performingUserID = uid.(string)
	} else {
		performingUserID = u.id
	}

	u.AddAudit(AuditUserCreated, performingUserID, nil, u.DTO())

	return u, nil
}

// ID returns the user's id.
func (u *User) ID() string {
	return u.id
}

// Email returns the User's email.
func (u *User) Email() string {
	return u.email
}

// NormalizedEmail returns the User's normalized email.
func (u *User) NormalizedEmail() string {
	return u.normalizedEmail
}

// Scopes returns the user's scopes.
func (u *User) Scopes() []*Scope {
	return u.scopes
}

// Update is used to update the user's core values, in a single function,
// by calling each other function. Update does not update the user's password.
func (u *User) Update(ctx context.Context, d *dto.UpdateUser, norm normalization.Normalizer) error {
	beforeUpdate := u.DTO()

	err := u.UpdateFirstname(d.Firstname)
	if err != nil {
		return err
	}

	err = u.UpdateLastname(d.Lastname)
	if err != nil {
		return err
	}

	err = u.UpdateEmail(d.Email, norm)
	if err != nil {
		return err
	}

	u.AddAudit(AuditUserUpdated, u.getPerformingUserID(ctx), beforeUpdate, u.DTO())

	return nil
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

// ChangePassword updates the user's password, ensuring the dto contains
// the user's current password.
func (u *User) ChangePassword(ctx context.Context, d *dto.ChangePassword, svc password.Service) result.Result {
	err := u.VerifyPassword(d.CurrentPassword, svc)
	if err != nil {
		msg := "Current password is invalid."
		return result.Failure(msg).WithStatusCode(http.StatusBadRequest)
	}

	err = u.setPassword(d.NewPassword, svc)
	if err != nil {
		return result.Failure(err).WithStatusCode(http.StatusBadRequest)
	}

	u.AddAudit(AuditUserPasswordChanged, u.getPerformingUserID(ctx), nil, nil)

	return result.Ok()
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

// ResetPassword is used to reset the user's password with a randomly generated string.
// A string will be generated then set as the user's password. The new password will be returned.
func (u *User) ResetPassword(ctx context.Context, scv password.Service) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#*&!?£$"
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 10)

	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	pwd := string(b)
	u.passwordHash = scv.Hash(pwd)

	u.AddAudit(AuditUserPasswordReset, u.getPerformingUserID(ctx), nil, nil)

	return pwd
}

func (u *User) getPerformingUserID(ctx context.Context) string {
	uid := ctx.Value(contextkey.ContextKey("user_id"))
	if uid != nil {
		return uid.(string)
	} else {
		return u.id
	}
}

func (u *User) VerifyPassword(password string, serv password.Service) error {
	ok := serv.Verify(password, u.passwordHash)
	if !ok {
		return fmt.Errorf("password is invalid")
	}

	return nil
}

// DataModel returns a datamodel object for the User.
func (u *User) DataModel() *datamodel.User {
	return &datamodel.User{
		ID:              u.id,
		Firstname:       u.firstname,
		Lastname:        u.lastname,
		Email:           u.email,
		NormalizedEmail: u.normalizedEmail,
		PasswordHash:    u.passwordHash,
	}
}

// UserFromDataModel returns a new instance of User populated with
// data from the given data-model object.
func UserFromDataModel(dm *datamodel.User, sdm []*datamodel.UserScope) *User {
	var scopes []*Scope
	if sdm != nil && len(sdm) > 0 {
		for _, s := range sdm {
			scopes = append(scopes, ScopeFromDataModel(s))
		}
	}

	return &User{
		id: dm.ID,
		firstname: dm.Firstname,
		lastname: dm.Lastname,
		email: dm.Email,
		normalizedEmail: dm.NormalizedEmail,
		passwordHash: dm.PasswordHash,
		scopes: scopes,
	}
}

// DTO returns a dto.User populated with the user' data.
func (u *User) DTO() *dto.User {
	return &dto.User{
		ID:              u.id,
		Firstname:       u.firstname,
		Lastname:        u.lastname,
		Email:           u.email,
		NormalizedEmail: u.normalizedEmail,
	}
}

// AddAudit adds an audit message/log to the user. The userID param
// is the id of the user who performed the action.
func (u *User) AddAudit(message, userID string, stateBefore, stateAfter *dto.User) {
	u.RaiseEvent(&event.AddUserAudit{
		Message:          message,
		Date:             time.Now().UTC(),
		UserID:           u.id,
		PerformingUserID: userID,
		Before: stateBefore,
		After: stateAfter,
	})
}