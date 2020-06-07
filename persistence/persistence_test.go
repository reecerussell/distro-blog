package persistence

import (
	"testing"

	"github.com/reecerussell/distro-blog/libraries/database"
)

func TestMySQLUserRepo(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("expected no panic, but got: %v", r)
		}
	}()

	db := database.NewMySQL("")
	_ = NewUserRepository(db)
}

func TestUnsupportedUserRepo(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected a panic")
		}
	}()

	_ = NewUserRepository("")
}
