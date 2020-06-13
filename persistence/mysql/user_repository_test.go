package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"os"
	"testing"

	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/normalization"
	"github.com/reecerussell/distro-blog/password"
)

var (
	testConnString = os.Getenv("CONN_STRING")
	testRepo       repository.UserRepository
)

func init() {
	db := database.NewMySQL(testConnString)
	testRepo = NewUserRepository(db)
}

func TestList(t *testing.T) {
	_, _, _, err := testRepo.List(context.Background()).Deconstruct()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetUser(t *testing.T) {
	_, id := seedUser("getUser@test.com")
	ctx := context.Background()
	// TODO: compare user values.
	success, _, _, err := testRepo.Get(ctx, id).Deconstruct()
	if !success {
		t.Errorf("unexpected failure: %v", err)
		return
	}

	t.Run("Not Found", func(t *testing.T) {
		success := testRepo.Get(ctx, uuid.New().String()).IsOk()
		if success{
			t.Errorf("expected to fail")
		}
	})
}

func TestAdd(t *testing.T) {
	u := buildUser("testAdd@test.com")
	ctx := context.Background()
	success, _, _, err := testRepo.Add(ctx, u).Deconstruct()
	if !success {
		t.Errorf("expected no error but got: %v", err)
	}
}

func TestAddWithExistingEmail(t *testing.T) {
	ctx := context.Background()
	testEmail := "addWithExistingEmail@test.com"

	executeHelper("INSERT INTO `users` (`id`,`first_name`,`last_name`,`email`,`normalized_email`,`password_hash`) VALUES (UUID(),?,?,?,?,?);",
		"John", "Doe", testEmail, normalization.New().Normalize(testEmail), "random string")

	// add duplicate user
	success := testRepo.Add(ctx, buildUser(testEmail)).IsOk()
	if success {
		t.Errorf("Inserted user successfully; this shouldn't work :/")
	}
}

func TestCountByEmail(t *testing.T) {
	ctx := context.Background()
	testEmail := "countByEmail@test.com"

	// count - assert 0
	success, _, count, err := testRepo.CountByEmail(ctx, buildUser(testEmail)).Deconstruct()
	if !success {
		t.Logf("Failed, expected no error but got: %v", err)
		return
	}

	if count.(int64) != 0 {
		t.Fail()
	}

	executeHelper("INSERT INTO `users` (`id`,`first_name`,`last_name`,`email`,`normalized_email`,`password_hash`) VALUES (UUID(),?,?,?,?,?);",
		"John", "Doe", testEmail, normalization.New().Normalize(testEmail), "random string")

	// count - assert 1
	success, _, count, err = testRepo.CountByEmail(ctx, buildUser(testEmail)).Deconstruct()
	if !success {
		return
	}

	if count.(int64) != 1 {
		t.Fail()
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

func seedUser(email string) (*model.User, string) {
	u := buildUser(email)
	dm := u.DataModel()
	executeHelper("CALL create_user(?,?,?,?,?,?);",
		dm.ID, dm.Lastname, dm.Lastname, dm.Email, dm.NormalizedEmail, dm.PasswordHash)

	return u, dm.ID
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

func countUsers() (c int64) {
	db, err := sql.Open("mysql", testConnString)
	if err != nil {
		panic(fmt.Errorf("open: %v", err))
	}

	err = db.QueryRow("select count(*) from users;").Scan(&c)
	if err != nil {
		panic(fmt.Errorf("query, scan: %v", err))
	}

	return
}
