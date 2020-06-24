package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/reecerussell/distro-blog/libraries/caching"
	"github.com/reecerussell/distro-blog/libraries/contextkey"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/logging"
	"github.com/reecerussell/distro-blog/libraries/storage"
	"github.com/reecerussell/distro-blog/persistence"
	"github.com/reecerussell/distro-blog/usecase"
	"os"
	"strings"
)

var (
	auth usecase.AuthUsecase
	store *storage.Service
	cache caching.Cache
)

type Resource struct {
	MethodARNSuffix string `yaml:"methodArnSuffix"`
	Scopes []string `yaml:"scopes"`
}

type Resources struct {
	Resources []Resource `yaml:"resources"`
}

func getAllowedScopes(methodArn string) []string {
	key := os.Getenv("SCOPES_STORAGE_KEY")
	methodKey := fmt.Sprintf("auth_scopes:%s", methodArn)
	bytes, ok := cache.Get(methodKey)

	var resources Resources

	if ok {
		var scopes []string
		err := json.Unmarshal(bytes, &scopes)
		if err == nil {
			return scopes
		}
	} else {
		bytes, ok = cache.Get(key)
		if ok {
			json.Unmarshal(bytes, &resources)
		} else {
			bytes, err := store.Get(key)
			if err != nil {
				err = fmt.Errorf("failed to get resources from store: %v", err)
				logging.Error(err)
				panic(err)
			}

			json.Unmarshal(bytes, &resources)
		}
	}

	var scopes []string

	for _ ,r := range resources.Resources {
		if !strings.HasSuffix(methodArn, r.MethodARNSuffix) {
			continue
		}

		scopes = r.Scopes
	}

	bytes, err := json.Marshal(scopes)
	if err == nil {
		err = cache.Set(methodKey, bytes)
		if err != nil {
			logging.Errorf("failed to set cache for method arn '%s': %v", methodArn, err)
		}
	}

	return scopes
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

	success, _, _, err := auth.VerifyWithScopes(ctx, []byte(parts[1]), getAllowedScopes(req.MethodArn)...).Deconstruct()
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
	db := database.NewMySQL(os.Getenv("CONN_STRING"))
	repo := persistence.NewUserRepository(db)
	auth = usecase.NewAuthUsecase(repo)

	var err error
	store, err = storage.New(os.Getenv("CONFIG_BUCKET_NAME"))
	if err != nil {
		err = fmt.Errorf("failed to init storage: %v", err)
		logging.Error(err)
		panic(err)
	}

	cache, err = caching.New(os.Getenv("CACHE_HOST"))
	if err != nil {
		err = fmt.Errorf("failed to init cache: %v", err)
		logging.Error(err)
		panic(err)
	}

	lambda.Start(handleAuthorization)
}