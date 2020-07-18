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

var settings usecase.SettingUsecase

func init() {
	db := database.NewMySQL(os.Getenv("CONN_STRING"))
	repo := persistence.NewSettingRepository(db)
	settings = usecase.NewSettingUsecase(repo)
}

func handleUpdateUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx = helper.PopulateContext(ctx, req)

	var us dto.Setting
	err := helper.ReadBody(req, &us)
	if err != nil {
		br := result.Failure(err).WithStatusCode(http.StatusBadRequest)
		return helper.Response(ctx, br, req), nil
	}

	res := settings.Update(ctx, &us)
	return helper.Response(ctx, res, req), nil
}

func main() {
	lambda.Start(handleUpdateUser)
}