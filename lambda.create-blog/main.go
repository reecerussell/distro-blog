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

var pages usecase.PageUsecase

func init() {
	db := database.NewMySQL(os.Getenv("CONN_STRING"))
	repo := persistence.NewPageRepository(db)
	pages = usecase.NewPageUsecase(repo, nil)
}

// handleCreate handles incoming API Gateway requests to create a new blog.
func handleCreate(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx = helper.PopulateContext(ctx, req)

	var d dto.CreatePage
	err := helper.ReadBody(req, &d)
	if err != nil {
		br := result.Failure(err).WithStatusCode(http.StatusBadRequest)
		return helper.Response(ctx, br, req), nil
	}

	res := pages.CreateBlog(ctx, &d)
	return helper.Response(ctx, res, req), nil
}

func main() {
	lambda.Start(handleCreate)
}
