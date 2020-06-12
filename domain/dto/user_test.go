package dto

import (
	"encoding/json"
	"testing"
)

var testUser = &User {
	ID: "273912",
	Firstname: "John",
	Lastname: "Doe",
	Email: "john@doe.com",
	NormalizedEmail: "JOHN@DOE.COM",
}

func TestMarshalUser(t *testing.T) {
	bytes, err := json.Marshal(testUser)
	if err != nil {
		t.Errorf("failed to marshal User: %v", err)
		return
	}

	var user User
	err = json.Unmarshal(bytes, &user)
	if err != nil {
		t.Errorf("failed to unmarshal User: %v", err)
		return
	}

	if user.ID != testUser.ID {
		t.Errorf("expected '%s' but got '%s'", testUser.ID, user.ID)
	}

	if user.Firstname != testUser.Firstname {
		t.Errorf("expected '%s' but got '%s'", testUser.Firstname, user.Firstname)
	}

	if user.Lastname != testUser.Lastname {
		t.Errorf("expected '%s' but got '%s'", testUser.Lastname, user.Lastname)
	}

	if user.Email != testUser.Email {
		t.Errorf("expected '%s' but got '%s'", testUser.Email, user.Email)
	}

	if user.NormalizedEmail != testUser.NormalizedEmail {
		t.Errorf("expected '%s' but got '%s'", testUser.NormalizedEmail, user.NormalizedEmail)
	}
}