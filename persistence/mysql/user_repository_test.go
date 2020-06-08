package mysql

import (
	"context"
	"database/sql"
	"fmt"
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

	t.Logf("Seeding the database with user: %s...", testEmail)
	executeHelper("INSERT INTO `users` (`id`,`first_name`,`last_name`,`email`,`normalized_email`,`password_hash`) VALUES (UUID(),?,?,?,?,?);",
		"John", "Doe", testEmail, normalization.New().Normalize(testEmail), "random string")

	// add duplicate user
	t.Logf("Attempting to create a user with a non-unique email...")
	u := buildUser(testEmail)
	success := testRepo.Add(ctx, u).IsOk()
	if success {
		t.Logf("\tsucceeded - should've failed!")
		t.Errorf("expected an error but got nil")
	} else {
		t.Logf("\t failed - expected!")
	}
}

func TestCountByEmail(t *testing.T) {
	ctx := context.Background()
	testEmail := "countByEmail@test.com"

	// count - assert 0
	t.Logf("Counting the number of users with email: %s...", testEmail)
	success, _, count, err := testRepo.CountByEmail(ctx, buildUser(testEmail)).Deconstruct()
	if !success {
		t.Logf("Failed, expected no error but got: %v", err)
		return
	}
	t.Logf("Expected 0, Actual: %d", count)

	if count.(int64) != 0 {
		t.Errorf("Expected 0, Actual: %v", count)
	}

	t.Logf("Seeding the database with user: %s...", testEmail)
	executeHelper("INSERT INTO `users` (`id`,`first_name`,`last_name`,`email`,`normalized_email`,`password_hash`) VALUES (UUID(),?,?,?,?,?);",
		"John", "Doe", testEmail, normalization.New().Normalize(testEmail), "random string")

	// count - assert 1
	t.Logf("Recounting users with email: %s...", testEmail)
	success, _, count, err = testRepo.CountByEmail(ctx, buildUser(testEmail)).Deconstruct()
	if !success {
		t.Logf("Failed, unexpected error: %v", err)
		return
	}
	t.Logf("Expected 1, Actual: %d", count)

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
