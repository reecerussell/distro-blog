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

var (
	testConnString = os.Getenv("CONN_STRING")
	users UserUsecase
)

func init() {
	db := database.NewMySQL(testConnString)
	repo := persistence.NewUserRepository(db)
	users = NewUserUsecase(repo)
}

func TestList(t *testing.T) {
	ctx := context.Background()
	res := users.List(ctx)
	if !res.IsOk() {
		t.Errorf("expected result to be successful")
	}
}

func TestCreate(t *testing.T) {
	defer executeHelper("delete from `users`;")

	db := database.NewMySQL(testConnString)
	repo := persistence.NewUserRepository(db)
	u := NewUserUsecase(repo)
	ctx := context.Background()
	cu := &dto.CreateUser{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "testCreate@test.com",
		Password:  "myTestPass-123",
	}

	success, _, _, err := u.Create(ctx, cu).Deconstruct()
	if !success {
		t.Errorf("expected no error but got: %v", err)
	}
}

func TestCreateWithInvalidData(t *testing.T) {
	executeHelper("DELETE FROM `users`;")
	defer executeHelper("delete from `users`;")

	db := database.NewMySQL(testConnString)
	repo := persistence.NewUserRepository(db)
	u := NewUserUsecase(repo)
	ctx := context.Background()
	cu := &dto.CreateUser{
		Firstname: "",
		Lastname:  "",
		Email:     "invalid email",
		Password:  "123",
	}

	success := u.Create(ctx, cu).IsOk()
	if success {
		t.Errorf("expected an error but got nil")
	}
}

func TestCreateWithExistingEmail(t *testing.T) {
	executeHelper("DELETE FROM `users`;")
	executeHelper("call create_user(?,?,?,?,?,?)", "0023823", "John", "Doe", "createWithExistingEmail@test.com", "CREATEWITHEXISTINGEMAIL@TEXT.CM", "password")
	defer executeHelper("delete from `users`;")

	db := database.NewMySQL(testConnString)
	repo := persistence.NewUserRepository(db)
	u := NewUserUsecase(repo)
	ctx := context.Background()
	cu := &dto.CreateUser{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "createWithExistingEmail@test.com",
		Password:  "myTestPass-123",
	}

	success := u.Create(ctx, cu).IsOk()
	if success {
		t.Errorf("expected an error but got nil")
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
