package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/helper"
	"github.com/reecerussell/distro-blog/persistence"
	"github.com/reecerussell/distro-blog/usecase"
	"os"
)

var users usecase.UserUsecase

func init() {
	db := database.NewMySQL(os.Getenv("CONN_STRING"))
	repo := persistence.NewUserRepository(db)
	users = usecase.NewUserUsecase(repo)
}

// handleDeleteUser is a handler function for Lambda, used to handle incoming proxy requests
// from an API gateway.
func handleDeleteUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx = helper.PopulateContext(ctx, req)
	id := req.PathParameters["id"]
	res := users.Delete(ctx, id)

	return helper.Response(ctx, res, req), nil
}

func main() {
	lambda.Start(handleDeleteUser)
}