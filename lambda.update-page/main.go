package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/reecerussell/aws-lambda-multipart-parser/parser"

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

	ir := persistence.NewImageRepository(db)
	itr := persistence.NewImageTypeRepository(db)
	media, err := usecase.NewMediaUsecase(ir, itr)
	if err != nil {
		panic(err)
	}

	pages = usecase.NewPageUsecase(repo, media)
}

// handleUpdate handles incoming API Gateway requests to update a page.
func handleUpdate(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx = helper.PopulateContext(ctx, req)

	var (
		d dto.UpdatePage
		imageData []byte = nil
	)

	if !isMultiPart(req) {
		err := helper.ReadBody(req, &d)
		if err != nil {
			br := result.Failure(err).WithStatusCode(http.StatusBadRequest)
			return helper.Response(ctx, br, req), nil
		}
	} else {
		data, err := parser.Parse(req)
		if err != nil {
			br := result.Failure(err).WithStatusCode(http.StatusBadRequest)
			return helper.Response(ctx, br, req), nil
		}

		var ok bool
		d.ID, ok = data.Get("id")
		if !ok {
			msg := "missing 'id' in request"
			br := result.Failure(msg).WithStatusCode(http.StatusBadRequest)
			return helper.Response(ctx, br, req), nil
		}

		d.Title, ok = data.Get("title")
		if !ok {
			msg := "missing 'title' in request"
			br := result.Failure(msg).WithStatusCode(http.StatusBadRequest)
			return helper.Response(ctx, br, req), nil
		}

		d.Description, ok = data.Get("description")
		if !ok {
			msg := "missing 'description' in request"
			br := result.Failure(msg).WithStatusCode(http.StatusBadRequest)
			return helper.Response(ctx, br, req), nil
		}

		content, ok := data.Get("content")
		if ok {
			d.Content = &content
		}

		file, ok := data.File("image")
		if ok {
			imageData = file.Content
		}
	}

	res := pages.Update(ctx, &d, imageData)
	return helper.Response(ctx, res, req), nil
}

func isMultiPart(req events.APIGatewayProxyRequest) bool {
	for k, v := range req.Headers {
		if strings.ToLower(k) == "content-type" {
			return strings.Contains(v, "multipart/form-data")
		}
	}

	return false
}

func main() {
	lambda.Start(handleUpdate)
}
