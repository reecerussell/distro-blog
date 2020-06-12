package dto

import (
	"encoding/json"
	"testing"
)

var testCreateUser = &CreateUser{
	Firstname: "john",
	Lastname:  "doe",
	Email:     "john@doe.com",
	Password:  "password123",
}

func TestMarshalCreateUser(t *testing.T) {
	bytes, err := json.Marshal(testCreateUser)
	if err != nil {
		t.Errorf("failed to marshal CreateUser: %v", err)
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