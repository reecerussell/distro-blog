package persistence

import (
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/persistence/mysql"
)

// New returns an instance of UserRepository for the given database type.
func New(db interface{}) repository.UserRepository {
	switch db.(type) {
	case *database.MySQL:
		return mysql.NewUserRepository(db.(*database.MySQL))
	default:
		panic("unsupported database type")
	}
}
