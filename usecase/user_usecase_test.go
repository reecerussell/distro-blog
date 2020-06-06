package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/persistence"
)

var testConnString = os.Getenv("CONN_STRING")

func TestCreate(t *testing.T) {
	defer executeHelper("delete from `users`;")

	db := database.NewMySQL(testConnString)
	repo := persistence.NewUserRepository(db)
	u := NewUserUsecase(repo)
	ctx := context.Background()
	cu := &dto.CreateUser{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john@doe.com",
		Password:  "myTestPass-123",
	}

	success, _, _, err := u.Create(ctx, cu).Deconstruct()
	if !success {
		t.Errorf("expected no error but got: %v", err)
	}
}

func executeHelper(query string, args ...interface{}) {
	db, err := sql.Open("mysql", testConnString)
	if err != nil {
		panic(fmt.Errorf("open: %v", err))
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		panic(fmt.Errorf("exec: %v", err))
	}
}
