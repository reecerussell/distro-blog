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
	bytes, err := json.Marshal(&testUserListItem)
	if err != nil {
		t.Errorf("failed to marshal UserListItem: %v", err)
		return
	}

	var item UserListItem
	err = json.Unmarshal(bytes, &item)
	if err != nil {
		t.Errorf("failed to unmarshal UserListItem: %v", err)
		return
	}

	if item.ID != testUserListItem.ID {
		t.Errorf("expected '%s' but got '%s'", testUserListItem.ID, item.ID)
	}

	if item.Name != testUserListItem.Name {
		t.Errorf("expected '%s' but got '%s'", testUserListItem.Name, item.Name)
	}

	if item.Email != testUserListItem.Email {
		t.Errorf("expected '%s' but got '%s'", testUserListItem.Email, item.Email)
	}
}