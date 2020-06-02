package dto_test

import (
	"encoding/json"
	"testing"

	"github.com/reecerussell/distro-blog/domain/dto"
)

var testCreateUser = &dto.CreateUser{
	Firstname: "john",
	Lastname:  "doe",
	Email:     "john@doe.com",
	Password:  "password123",
}

func TestMarshal(t *testing.T) {
	bytes, err := json.Marshal(testCreateUser)
	if err != nil {
		t.Errorf("failed to marshal dto.CreateUser: %v", err)
		return
	}

	var data map[string]string
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		t.Errorf("failed to unmarshal dto.CreateUser to map[string]string: %v", err)
		return
	}

	if v := data["firstname"]; v != testCreateUser.Firstname {
		t.Errorf("expected '%s' but got '%s'", testCreateUser.Firstname, v)
	}

	if v := data["lastname"]; v != testCreateUser.Lastname {
		t.Errorf("expected '%s' but got '%s'", testCreateUser.Lastname, v)
	}

	if v := data["email"]; v != testCreateUser.Email {
		t.Errorf("expected '%s' but got '%s'", testCreateUser.Email, v)
	}

	if v := data["password"]; v != testCreateUser.Password {
		t.Errorf("expected '%s' but got '%s'", testCreateUser.Password, v)
	}
}

func TestUnmarshal(t *testing.T) {
	bytes, err := json.Marshal(testCreateUser)
	if err != nil {
		t.Errorf("failed to marshal dto.CreateUser: %v", err)
		return
	}

	var createUser dto.CreateUser
	err = json.Unmarshal(bytes, &createUser)
	if err != nil {
		t.Errorf("failed to unmarshal dto.CreateUser: %v", err)
		return
	}

	if createUser.Firstname != testCreateUser.Firstname {
		t.Errorf("expected '%s' but got '%s'", testCreateUser.Firstname, createUser.Firstname)
	}

	if createUser.Lastname != testCreateUser.Lastname {
		t.Errorf("expected '%s' but got '%s'", testCreateUser.Lastname, createUser.Lastname)
	}

	if createUser.Email != testCreateUser.Email {
		t.Errorf("expected '%s' but got '%s'", testCreateUser.Email, createUser.Email)
	}

	if createUser.Password != testCreateUser.Password {
		t.Errorf("expected '%s' but got '%s'", testCreateUser.Password, createUser.Password)
	}
}
