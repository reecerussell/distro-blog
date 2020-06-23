package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"net/http"
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
	testConnStringEmptySchema = os.Getenv("CONN_STRING_EMPTY_SCHEMA")
	testConnStringDeformed = os.Getenv("CONN_STRING_DEFORMED")
	testRepo       repository.UserRepository
	testRepoEmptySchema repository.UserRepository
	testRepoDeformed repository.UserRepository
)

func init() {
	db := database.NewMySQL(testConnString)
	testRepo = NewUserRepository(db)
	db = database.NewMySQL(testConnStringEmptySchema)
	testRepoEmptySchema = NewUserRepository(db)
	db = database.NewMySQL(testConnStringDeformed)
	testRepoDeformed = NewUserRepository(db)
}

func TestList(t *testing.T) {
	_, _, _, err := testRepo.List(context.Background()).Deconstruct()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUserRepository_UserReader(t *testing.T) {
	seedUser("userRepository@userReader.test")

	success, _, _, err := testRepo.List(context.Background()).Deconstruct()
	if !success {
		t.Errorf("unexpected error: %v", err)
	}

	t.Run("Missing Schema", func(t *testing.T) {
		res := testRepoEmptySchema.List(context.Background())
		if res.IsOk(){
			t.Errorf("expected to fail")
		}
	})

	t.Run("Deformed View", func(t *testing.T) {
		res := testRepoDeformed.List(context.Background())
		if res.IsOk(){
			t.Errorf("expected to fail")
		}
	})
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

	t.Run("Missing Table", func(t *testing.T) {
		res := testRepoEmptySchema.Get(ctx, "some id")
		if res.IsOk(){
			t.Errorf("expecte to fail")
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

func TestUserRepository_CountByEmail(t *testing.T) {
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

	t.Run("Missing Schema", func(t *testing.T) {
		res := testRepoEmptySchema.CountByEmail(ctx, buildUser(testEmail))
		if res.IsOk() {
			t.Errorf("expected to fail")
		}
	})
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

func TestUserRepository_Update(t *testing.T) {
	u, _ := seedUser("userRepository@update.test")
	ud := &dto.UpdateUser{
		Firstname: "UpdateRepository",
		Lastname: "Test",
		Email: u.Email(),
	}
	_ = u.Update(ud, normalization.New())
	ctx := context.Background()

	success, _, _, err := testRepo.Update(ctx, u).Deconstruct()
	if !success {
		t.Errorf("unexpected error: %v", err)
		return
	}

	// TODO: check user values.

	t.Run("Non Existent User", func(t *testing.T) {
		u := buildUser("nonExistentUser@update.test")
		success, status, _, _ := testRepo.Update(ctx, u).Deconstruct()
		if success {
			t.Errorf("unexpected success")
		}

		if status != http.StatusNotFound {
			t.Errorf("expected status code %d, but got %d", http.StatusNotFound, status)
		}
	})

	t.Run("Missing Stored Procedure", func(t *testing.T) {
		u := buildUser("missingSproc@update.test")
		ctx := context.Background()

		res := testRepoEmptySchema.Update(ctx, u)
		if res.IsOk() {
			t.Errorf("expected an error")
		}
	})
}

func TestUserRepository_Delete(t *testing.T) {
	testEmail := "userRepository@delete.test"
	_, id := seedUser(testEmail)
	ctx := context.Background()

	success, _, _, err := testRepo.Delete(ctx, id).Deconstruct()
	if !success {
		t.Errorf("unexpected error: %v", err)
	}

	t.Run("Invalid Schema", func(t *testing.T) {
		res := testRepoEmptySchema.Delete(ctx, id)
		if res.IsOk(){
			t.Errorf("expected an error")
		}
	})
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
