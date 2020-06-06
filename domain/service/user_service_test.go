package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/normalization"
	"github.com/reecerussell/distro-blog/password"
	"github.com/reecerussell/distro-blog/persistence/mysql"
)

var (
	testConnString = os.Getenv("CONN_STRING")
)

func TestEnsureEmailIsUnique(t *testing.T) {
	defer executeHelper("delete from `users`;")

	db := database.NewMySQL(testConnString)
	repo := mysql.NewUserRepository(db)
	serv := NewUserService(repo)
	ctx := context.Background()

	u := buildUser()
	success, _, _, err := serv.EnsureEmailIsUnique(ctx, u).Deconstruct()
	if !success {
		t.Errorf("expected no error but got: %v", err)
	}
}

func TestEnsureEmailIsUniqueWithNonUnique(t *testing.T) {
	defer executeHelper("delete from `users`;")

	db := database.NewMySQL(testConnString)
	repo := mysql.NewUserRepository(db)
	serv := NewUserService(repo)
	ctx := context.Background()

	// seed user
	u := buildUser()
	success, _, _, err := repo.Add(ctx, u).Deconstruct()
	if !success {
		t.Errorf("expected no error but got: %v", err)
	}

	// check with existing user
	success, _, _, err = serv.EnsureEmailIsUnique(ctx, u).Deconstruct()
	if !success {
		t.Errorf("expected no error but got: %v", err)
	}

	// check new user
	u = buildUser()
	success, _, _, _ = serv.EnsureEmailIsUnique(ctx, u).Deconstruct()
	if success {
		t.Errorf("expected an error but got nil")
	}
}

func buildUser() *model.User {
	cu := &dto.CreateUser{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john@doe.com",
		Password:  "MyTestPassword123",
	}
	norm := normalization.New()
	pwdServ := password.New()
	u, err := model.NewUser(cu, pwdServ, norm)
	if err != nil {
		panic(err)
	}

	return u
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
