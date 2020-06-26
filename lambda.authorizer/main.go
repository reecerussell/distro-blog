package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/yaml.v2"

	"github.com/reecerussell/distro-blog/libraries/contextkey"
	"github.com/reecerussell/distro-blog/libraries/logging"
	"github.com/reecerussell/distro-blog/libraries/storage"
	"github.com/reecerussell/distro-blog/usecase"
)

var (
	auth usecase.AuthUsecase
	store *storage.Service
	config *Config
)

func init(){
	auth = usecase.NewAuthUsecase(nil)

	var err error
	store, err = storage.New(os.Getenv("CONFIG_BUCKET_NAME"))
	if err != nil {
		err = fmt.Errorf("failed to init storage: %v", err)
		logging.Error(err)
		panic(err)
	}

	config, err = loadConfig(os.Getenv("CONFIG_BUCKET_KEY"))
	if err != nil {
		logging.Error(err)
		panic(err)
	}
}

func buildResources(methodArn string, scopes []string) []string {
	resourceMap := make(map[string]bool)
	parts := strings.Split(methodArn, "/")
	baseArn := strings.Join(parts[:2], "/")

	var resources []string

	for _, scope := range scopes {
		methods, ok := config.ScopePolicies[scope]
		if !ok {
			continue
		}

		for _, m := range methods {
			_, ok = resourceMap[baseArn + m]
			if ok {
				continue
			}

			resourceMap[baseArn + m] = true
			resources = append(resources, baseArn+m)
		}
	}

	return resources
}

func findAllowedScopes(methodArn string) []string {
	allowedMap := make(map[string]bool)

	var allowed []string

	for suf, scps := range config.Scopes {
		suf = strings.ReplaceAll(suf, "*", "(.+)")
		arn := methodArn[strings.Index(methodArn, "/"):]

		re := regexp.MustCompile(fmt.Sprintf("%s$", suf))
		if !re.MatchString(arn) {
			continue
		}

		for _, s := range scps {
			_, ok := allowedMap[s]
			if ok {
				continue
			}

			allowedMap[s] = true
			allowed = append(allowed, s)
		}

		break
	}

	return allowed
}

func generatePolicy(effect, methodArn string, scopes []string) events.APIGatewayCustomAuthorizerResponse {
		return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "user",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action: []string{"execute-api:Invoke"},
					Effect: effect,
					Resource: buildResources(methodArn, scopes),
				},
			},
		},
	}
}

func handleAuthorization(ctx context.Context, req events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	ctx = context.WithValue(ctx, contextkey.ContextKey("JWT_KEY_ID"), os.Getenv("JWT_KEY_ID"))

	scopes := findAllowedScopes(req.MethodArn)
	parts := strings.Split(req.AuthorizationToken, " ")
	if len(parts) < 2 || parts[0] != "Bearer" {
		return generatePolicy("Deny", req.MethodArn, scopes), errors.New("Unauthorized")
	}

	success, status, _, err := auth.VerifyWithScopes(ctx, []byte(parts[1]), scopes...).Deconstruct()
	if !success {
		pol := generatePolicy("Deny", req.MethodArn, scopes)

		if status == http.StatusForbidden {
			err = errors.New("Forbidden")
		} else {
			err = errors.New("Unauthorized")
		}

		return pol, err
	}

	return generatePolicy("Allow", req.MethodArn, scopes), nil
}

type Config struct {
	Scopes map[string][]string `yaml:"scopes"`
	ScopePolicies map[string][]string `yaml:"scope_policies"`
}

func loadConfig(configKey string) (*Config, error) {
	data, err := store.Get(configKey)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config")
	}

	return &config, nil
}

func main() {
	lambda.Start(handleAuthorization)
}