package main

import (
	"context"
	"github.com/reecerussell/distro-blog/libraries/logging"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/helper"
	"github.com/reecerussell/distro-blog/persistence/mysql"
	"github.com/reecerussell/distro-blog/usecase"
)

var users usecase.UserUsecase

func init() {
	db := database.NewMySQL(os.Getenv("CONN_STRING"))
	repo := mysql.NewUserRepository(db)
	users = usecase.NewUserUsecase(repo)
}

func handleGetUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx = helper.PopulateContext(ctx, req)
	id := req.PathParameters["id"]
	expand := req.MultiValueQueryStringParameters["expand"]

	logging.Debugf("Expand: %s\n", strings.Join(expand, ", "))

	res := users.Get(ctx, id, expand...)

	return helper.Response(ctx, res, req), nil
}

func main() {
	lambda.Start(handleGetUser)
}