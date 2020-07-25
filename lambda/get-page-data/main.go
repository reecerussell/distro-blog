package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/helper"
	"github.com/reecerussell/distro-blog/libraries/result"
)

var (
	db *database.MySQL
)

func init() {
	db = database.NewMySQL(os.Getenv("CONN_STRING"))
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx = helper.PopulateContext(ctx, req)
	res := getPageData(ctx, req.PathParameters["url"])
	return helper.Response(ctx, res, req), nil
}

func main() {
	lambda.Start(handler)
}

func getPageData(ctx context.Context, url string) result.Result {
	url = strings.ToLower(url)
	if strings.HasPrefix(url, "blog-") {
		url = "blog/" + url[5:]
	}

	const query string = "CALL `get_page_data_by_url`(?);"
	data, err := db.Read(ctx, query, pageReader, url)
	if err != nil {
		return result.Failure(err)
	}

	if data == nil {
		return result.Failure("page not found").WithStatusCode(http.StatusNotFound)
	}

	return result.Ok().WithValue(data)
}

type PageData struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Content *string `json:"content"`
	IsBlog bool `json:"isBlog"`
	ImageID *string `json:"imageId"`
	SEO SEOData `json:"seo"`
}

type SEOData struct {
	Title string `json:"title"`
	Description string `json:"description"`
	SiteName string `json:"siteName"`
	Index bool `json:"index"`
	Follow bool `json:"follow"`
}

func pageReader(s database.ScannerFunc) (interface{}, error) {
	var data PageData
	var content sql.NullString
	var imageID sql.NullString
	err := s(
		&data.ID,
		&data.Title,
		&data.Description,
		&content,
		&data.IsBlog,
		&imageID,
		&data.SEO.Title,
		&data.SEO.Description,
		&data.SEO.SiteName,
		&data.SEO.Index,
		&data.SEO.Follow,
	)
	if err != nil {
		return nil, err
	}

	if content.Valid {
		data.Content = &content.String
	}

	if imageID.Valid {
		data.ImageID = &imageID.String
	}

	return data, nil
}