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

func TestMySQLPageRepo(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("expected no panic, but got: %v", r)
		}
	}()

	db := database.NewMySQL("")
	_ = NewPageRepository(db)
}

func TestUnsupportedPageRepo(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected a panic")
		}
	}()

	_ = NewPageRepository("")
}

func TestNewImageRepository(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("expected no panic, but got: %v", r)
		}
	}()

	db := database.NewMySQL("")
	_ = NewImageRepository(db)

	t.Run("Unsupported Database Type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected a panic")
			}
		}()

		_ = NewImageRepository("")
	})
}

func TestNewImageTypeRepository(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("expected no panic, but got: %v", r)
		}
	}()

	db := database.NewMySQL("")
	_ = NewImageTypeRepository(db)

	t.Run("Unsupported Database Type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected a panic")
			}
		}()

		_ = NewImageTypeRepository("")
	})
}

func TestNewSettingRepository(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("expected no panic, but got: %v", r)
		}
	}()

	db := database.NewMySQL("")
	_ = NewSettingRepository(db)

	t.Run("Unsupported Database Type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected a panic")
			}
		}()

		_ = NewSettingRepository("")
	})
}