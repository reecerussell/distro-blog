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