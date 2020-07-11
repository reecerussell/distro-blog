package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/persistence"
	"github.com/reecerussell/distro-blog/usecase"
)

var media usecase.MediaUsecase

func init() {
	db := database.NewMySQL(os.Getenv("CONN_STRING"))
	ir := persistence.NewImageRepository(db)
	itr := persistence.NewImageTypeRepository(db)

	var err error
	media, err = usecase.NewMediaUsecase(ir, itr)
	if err != nil {
		panic(err)
	}
}

func handleDownload(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return media.DownloadForLambda(ctx, req.PathParameters["id"]), nil
}

func main() {
	lambda.Start(handleDownload)
}