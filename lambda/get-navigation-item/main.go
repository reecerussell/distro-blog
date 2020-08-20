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

var navigation usecase.NavigationUsecase

func init() {
	db := database.NewMySQL(os.Getenv("CONN_STRING"))
	repo := persistence.NewNavigationRepository(db)
	navigation = usecase.NewNavigationUsecase(repo)
}

// handleGet handles incoming API Gateway requests to get a navigation item.
func handleGet(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx = helper.PopulateContext(ctx, req)
	res := navigation.GetItem(ctx, req.PathParameters["id"])
	return helper.Response(ctx, res, req), nil
}

func main() {
	lambda.Start(handleGet)
}
