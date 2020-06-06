package mysql

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
)

var (
	testConnString = os.Getenv("CONN_STRING")
)

func TestAdd(t *testing.T) {
	defer executeHelper("delete from `users`;")

	db := database.NewMySQL(testConnString)
	repo := NewUserRepository(db)

	u := buildUser()
	ctx := context.Background()
	success, _, _, err := repo.Add(ctx, u).Deconstruct()
	if !success {
		t.Errorf("expected no error but got: %v", err)
	}
}

func TestAddWithExistingEmail(t *testing.T) {
	defer executeHelper("delete from `users`;")

	db := database.NewMySQL(testConnString)
	repo := NewUserRepository(db)
	ctx := context.Background()
	u := buildUser()

	// seed user
	success, _, _, err := repo.Add(ctx, u).Deconstruct()
	if !success {
		t.Errorf("expected no error but got: %v", err)
	}

	// add duplicate user
	success, _, _, err = repo.Add(ctx, u).Deconstruct()
	if success {
		t.Errorf("expected an error but got nil")
	}
}

func TestCountByEmail(t *testing.T) {
	defer executeHelper("delete from `users`;")

	db := database.NewMySQL(testConnString)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// count - assert 0
	u := buildUser()
	success, _, count, err := repo.CountByEmail(ctx, u).Deconstruct()
	if !success {
		t.Errorf("expected no error but got: %v", err)
	}

	if count.(int64) != 0 {
		t.Errorf("expected 0 but got: %v", count)
	}

	// add user
	success, _, _, err = repo.Add(ctx, u).Deconstruct()
	if !success {
		t.Errorf("expected no error but got: %v", err)
	}

	// count - assert 1
	u = buildUser()
	success, _, count, err = repo.CountByEmail(ctx, u).Deconstruct()
	if !success {
		t.Errorf("expected no error but got: %v", err)
	}

	if count.(int64) != 1 {
		t.Errorf("expected 1 but got: %v", count)
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
