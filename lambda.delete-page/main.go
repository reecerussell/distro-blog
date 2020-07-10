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

	ir := persistence.NewImageRepository(db)
	itr := persistence.NewImageTypeRepository(db)
	media, err := usecase.NewMediaUsecase(ir, itr)
	if err != nil {
		panic(err)
	}

	pages = usecase.NewPageUsecase(repo, media)
}

// handleDelete handles incoming API Gateway requests to delete pages.
func handleDelete(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx = helper.PopulateContext(ctx, req)
	res := pages.Delete(ctx, req.PathParameters["id"])
	return helper.Response(ctx, res, req), nil
}

func main() {
	lambda.Start(handleDelete)
}

