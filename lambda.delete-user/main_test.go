package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/reecerussell/distro-blog/libraries/normalization"
	"net/http"
	"testing"
	"database/sql"
	"os"
)

var(
	testConnString = os.Getenv("CONN_STRING")
)

func TestHandler(t *testing.T) {
	id := seedUser("deleteUser@lambda.test")

	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"id": id,
		},
	}
	resp, _ := handleDeleteUser(context.Background(), req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected a status code of %d but got %d", http.StatusOK, resp.StatusCode)
	}
}

func seedUser(email string) string {
	db, err := sql.Open("mysql", testConnString)
	if err != nil {
		panic(fmt.Errorf("open: %v", err))
	}

	id := uuid.New().String()
	_, err = db.Exec("CALL `create_user`(?,?,?,?,?,?);",
		id,
		"John",
		"Doe",
		email,
		normalization.New().Normalize(email),
		"763tegdjwhd")
	if err != nil {
		panic(fmt.Errorf("exec: %v", err))
	}

	return id
}