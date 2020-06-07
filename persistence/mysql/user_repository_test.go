package mysql

import (
	"context"
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

	t.Logf("Creating initial user: %s...", testEmail)
	success, _, _, err := testRepo.Add(ctx, buildUser(testEmail)).Deconstruct()
	if !success {
		t.Logf("\t failed.\n\t%v", err)
	} else {
		t.Logf("\t done.")
	}

	// add duplicate user
	t.Logf("Attempting to create a user with a non-unique email...")
	u := buildUser(testEmail)
	success = testRepo.Add(ctx, u).IsOk()
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
	u := buildUser(testEmail)
	success, _, count, err := testRepo.CountByEmail(ctx, u).Deconstruct()
	if !success {
		t.Logf("\tfailed.\n\t%v", err)
		t.Errorf("expected no error but got: %v", err)
		return
	}
	t.Logf("\t%v - expected 0\n", count)

	if count.(int64) != 0 {
		t.Errorf("expected 0 but got: %v", count)
	}

	// add user
	t.Logf("Creating initial user with email: %s...", testEmail)
	success, _, _, err = testRepo.Add(ctx, buildUser(testEmail)).Deconstruct()
	if !success {
		t.Logf("\tfailed - should've worked\n\t%v", err)
		t.Errorf("expected to be able to insert user.")
		return
	}

	// count - assert 1
	t.Logf("Recounting users with email: %s...", testEmail)
	u = buildUser(testEmail)
	success, _, count, err = testRepo.CountByEmail(ctx, u).Deconstruct()
	if !success {
		t.Logf("\tfailed - should've worked\n\t%v", err)
		t.Errorf("expected no error but got: %v", err)
		return
	} else {
		t.Logf("\t%v - expected 1\n", count)
	}

	if count.(int64) != 1 {
		t.Errorf("expected 1 but got: %v", count)
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
