package persistence

import (
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/persistence/mysql"
)

// NewUserRepository returns an instance of UserRepository for the given database type.
func NewUserRepository(db interface{}) repository.UserRepository {
	switch db.(type) {
	case *database.MySQL:
		return mysql.NewUserRepository(db.(*database.MySQL))
	default:
		panic("unsupported database type")
	}
}

// NewPageRepository returns an instance of PageRepository for the given database type.
func NewPageRepository(db interface{}) repository.PageRepository {
	switch db.(type) {
	case *database.MySQL:
		return mysql.NewPageRepository(db.(*database.MySQL))
	default:
		panic("unsupported database type")
	}
}

// NewImageRepository returns and instance of ImageRepository for the given database type.
func NewImageRepository(db interface{}) repository.ImageRepository {
	switch db.(type) {
	case *database.MySQL:
		return mysql.NewImageRepository(db.(*database.MySQL))
	default:
		panic("unsupported database type")
	}
}

// NewImageTypeRepository returns and instance of ImageTypeRepository for the given database type.
func NewImageTypeRepository(db interface{}) repository.ImageTypeRepository {
	switch db.(type) {
	case *database.MySQL:
		return mysql.NewImageTypeRepository(db.(*database.MySQL))
	default:
		panic("unsupported database type")
	}
}

// NewSettingRepository returns and instance of SettingRepository for the given database type.
func NewSettingRepository(db interface{}) repository.SettingRepository {
	switch db.(type) {
	case *database.MySQL:
		return mysql.NewSettingRepository(db.(*database.MySQL))
	default:
		panic("unsupported database type")
	}
}

// NewNavigationRepository returns and instance of NavigationRepository for the given database type.
func NewNavigationRepository(db interface{}) repository.NavigationRepository {
	switch db.(type) {
	case *database.MySQL:
		return mysql.NewNavigationRepository(db.(*database.MySQL))
	default:
		panic("unsupported database type")
	}
}