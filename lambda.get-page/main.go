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

var pages usecase.PageUsecase

func init() {
	db := database.NewMySQL(os.Getenv("CONN_STRING"))
	repo := persistence.NewPageRepository(db)
	pages = usecase.NewPageUsecase(repo, nil)
}

// handleGet handles incoming API Gateway requests to get a page.
func handleGet(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx = helper.PopulateContext(ctx, req)
	res := pages.Get(ctx, req.PathParameters["id"], req.MultiValueQueryStringParameters["expand"]...)
	return helper.Response(ctx, res, req), nil
}

func main() {
	lambda.Start(handleGet)
}
