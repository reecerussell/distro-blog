package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/helper"
	"github.com/reecerussell/distro-blog/libraries/logging"
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
		var ok bool
		imageData, ok = readFormData(req, &d)
		if !ok {
			br := result.Failure("invalid request body").WithStatusCode(http.StatusBadRequest)
			return helper.Response(ctx, br, req), nil
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

func readFormData(req events.APIGatewayProxyRequest, d *dto.UpdatePage) ([]byte, bool) {
	body := req.Body
	if req.IsBase64Encoded {
		bytes, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			logging.Errorf("failed to read base64 body: %v\n", err)
			return nil, false
		}

		body = string(bytes)
	}

	r := &http.Request{
		Header: make(map[string][]string),
	}
	for k, v := range req.Headers {
		r.Header.Set(k, v)
	}

	r.Body = ioutil.NopCloser(strings.NewReader(body))
	err := r.ParseMultipartForm(256 << 20)
	if err != nil {
		logging.Error(err)
		return nil, false
	}

	d.ID = r.FormValue("id")
	d.Title = r.FormValue("title")
	d.Description = r.FormValue("description")
	d.URL = r.FormValue("url")

	var seo dto.SEO
	if v := r.FormValue("seoTitle"); v != "" {
		seo.Title = &v
	}

	if v := r.FormValue("seoDescription"); v != "" {
		seo.Description = &v
	}

	if v := r.FormValue("seoIndex"); strings.ToLower(v) == "true" {
		seo.Index = true
	}

	if v := r.FormValue("seoFollow"); strings.ToLower(v) == "true" {
		seo.Follow = true
	}

	d.SEO = &seo

	content := r.FormValue("content")
	if content != "" {
		d.Content = &content
	}

	var img []byte

	file, hdr, err := r.FormFile("image")
	if err != nil {
		logging.Error(err)
	} else {
		defer file.Close()
		buf := make([]byte, hdr.Size)
		_, err = file.Read(buf)
		if err != nil {
			logging.Error(err)
		} else {
			img = buf
		}
	}

	return img, true
}

func getFormValue(form *multipart.Form, name string) (string, error) {
	values, ok := form.Value[name]
	if !ok || len(values) < 1 {
		return "", fmt.Errorf("missing the '%s' field", name)
	}

	return values[0], nil
}

func main() {
	lambda.Start(handleUpdate)
}
