package main

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/reecerussell/distro-blog/libraries/contextkey"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/persistence"
	"github.com/reecerussell/distro-blog/usecase"
	"os"
	"strings"
)

var (
	auth usecase.AuthUsecase
	scopes map[string][]string
)

func init() {
	db := database.NewMySQL(os.Getenv("CONN_STRING"))
	repo := persistence.NewUserRepository(db)
	auth = usecase.NewAuthUsecase(repo)
}

func generatePolicy(effect, methodArn string) events.APIGatewayCustomAuthorizerResponse {
	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "user",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action: []string{"execute-api:Invoke"},
					Effect: effect,
					Resource: []string{methodArn},
				},
			},
		},
	}
}

func handleAuthorization(ctx context.Context, req events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	ctx = context.WithValue(ctx, contextkey.ContextKey("JWT_KEY_ID"), os.Getenv("JWT_KEY_ID"))

	parts := strings.Split(req.AuthorizationToken, " ")
	if len(parts) < 2 || parts[0] != "Bearer" {
		return generatePolicy("Deny", req.MethodArn), errors.New("Unauthorized")
	}

	success, _, _, err := auth.VerifyWithScopes(ctx, []byte(parts[1])).Deconstruct()
	if !success {
		pol := generatePolicy("Deny", req.MethodArn)
		pol.Context = map[string]interface{}{
			"error": err.Error(),
		}
		return pol, errors.New("Unauthorized")
	}

	return generatePolicy("Allow", req.MethodArn), nil
}

func main() {
	lambda.Start(handleAuthorization)
}