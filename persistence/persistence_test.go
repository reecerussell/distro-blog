package persistence

import (
	"testing"

	"github.com/reecerussell/distro-blog/libraries/database"
)

func TestMySQL(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("expected no panic, but got: %v", r)
		}
	}()

	db := database.NewMySQL("")
	_ = New(db)
}

func TestUnsupported(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected a panic")
		}
	}()

	_ = New("")
}
