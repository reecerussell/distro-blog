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
	testServ       *UserService
)

func init() {
	db := database.NewMySQL(testConnString)
	repo := mysql.NewUserRepository(db)
	testServ = NewUserService(repo)
}

func TestEnsureEmailIsUnique(t *testing.T) {
	ctx := context.Background()

	u := buildUser("ensureEmailIsUnique@test.com")
	success, _, _, err := testServ.EnsureEmailIsUnique(ctx, u).Deconstruct()
	if !success {
		t.Errorf("expected no error but got: %v", err)
	}
}

func TestEnsureEmailIsUniqueWithNonUnique(t *testing.T) {
	ctx := context.Background()
	testEmail := "ensureEmailIsUniqueWithNonUnique@test.com"

	createUser(testEmail)

	u := buildUser(testEmail)
	success := testServ.EnsureEmailIsUnique(ctx, u).IsOk()
	if success {
		t.Errorf("expected an error but got nil")
	}
}

func buildUser(email string) *model.User {
	cu := &dto.CreateUser{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     email,
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

func createUser(email string) {
	norm := normalization.New()
	executeHelper("CALL `create_user`(UUID(),?,?,?,?,?);", "John", "Doe", email, norm.Normalize(email), "e3ije")
}

func executeHelper(query string, args ...interface{}) {
	fmt.Printf("Executing: %s", query)
	db, err := sql.Open("mysql", testConnString)
	if err != nil {
		panic(fmt.Errorf("open: %v", err))
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		panic(fmt.Errorf("exec: %v", err))
	}

	fmt.Printf("\t done.\n")
}
