package dto

import (
	"encoding/json"
	"testing"
)

var testUserListItem = &UserListItem{
	ID: "01283",
	Name: "John Doe",
	Email: "john@doe.com",
}

func TestMarshalUserListItem(t *testing.T) {
	bytes, err := json.Marshal(testCreateUser)
	if err != nil {
		t.Errorf("failed to marshal UserListItem: %v", err)
		return
	}

	var data map[string]string
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		t.Errorf("failed to unmarshal UserListItem: %v", err)
		return
	}

	if v := data["id"]; v != testUserListItem.ID {
		t.Errorf("expected '%s' but got '%s'", testUserListItem.ID, v)
	}

	if v := data["name"]; v != testUserListItem.Name {
		t.Errorf("expected '%s' but got '%s'", testUserListItem.Name, v)
	}

	if v := data["email"]; v != testUserListItem.Email {
		t.Errorf("expected '%s' but got '%s'", testUserListItem.Email, v)
	}
}