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

// handleList is a Lambda handler function used to handle incoming
// APIGateway proxy requests to gather a list of navigation items.
func handleList(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx = helper.PopulateContext(ctx, req)
	res := navigation.GetItems(ctx)
	return helper.Response(ctx, res, req), nil
}

func main() {
	lambda.Start(handleList)
}
