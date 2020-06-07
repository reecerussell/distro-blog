package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

var users usecase.UserUsecase

func init() {
	db := database.NewMySQL(os.Getenv("CONN_STRING"))
	repo := persistence.NewUserRepository(db)
	users = usecase.NewUserUsecase(repo)
}

// handleCreateUser is a Lambda handler function used to handle incoming APIGateway
// proxy requests, to create users.
func handleCreateUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx = helper.PopulateContext(ctx, req)

	var cu dto.CreateUser
	data, _ := base64.StdEncoding.DecodeString(req.Body)
	err := json.Unmarshal(data, &cu)
	if err != nil {
		msg := fmt.Sprintf("Failed to read request body, as it was in an unsupported format: %v", err)
		br := result.Failure(msg).WithStatusCode(http.StatusBadRequest)
		return helper.Response(ctx, br, req), nil
	}

	res := users.Create(ctx, &cu)
	return helper.Response(ctx, res, req), nil
}

func main() {
	lambda.Start(handleCreateUser)
}
