package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/helper"
	"github.com/reecerussell/distro-blog/persistence"
	"github.com/reecerussell/distro-blog/usecase"
)

var settings usecase.SettingUsecase

func init() {
	db := database.NewMySQL(os.Getenv("CONN_STRING"))
	repo := persistence.NewSettingRepository(db)
	settings = usecase.NewSettingUsecase(repo)
}

// handleGet handles incoming API Gateway requests to get a setting.
func handleGet(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx = helper.PopulateContext(ctx, req)
	res := settings.Get(ctx, req.PathParameters["key"])
	return helper.Response(ctx, res, req), nil
}

func main() {
	lambda.Start(handleGet)
}
