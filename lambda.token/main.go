package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/helper"
	"github.com/reecerussell/distro-blog/libraries/result"
	"github.com/reecerussell/distro-blog/persistence"
	"github.com/reecerussell/distro-blog/usecase"
)

var (
	auth usecase.AuthUsecase
)

func init() {
	db := database.NewMySQL(os.Getenv("CONN_STRING"))
	repo := persistence.NewUserRepository(db)
	auth = usecase.NewAuthUsecase(repo)
}

func handleToken(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx = helper.PopulateContext(ctx, req)

	var cred dto.UserCredential
	err := helper.ReadBody(req, &cred)
	if err != nil {
		br := result.Failure(err).WithStatusCode(http.StatusBadRequest)
		return helper.Response(ctx, br, req), nil
	}

	res := auth.Token(ctx, &cred)
	return helper.Response(ctx, res, req), nil
}

func main() {
	lambda.Start(handleToken)
}