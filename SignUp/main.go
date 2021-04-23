package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/igvaquero18/smarthome/controller"
	"github.com/igvaquero18/smarthome/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	awsRegionEnv         = "SMARTHOME_AWS_REGION"
	verboseEnv           = "SMARTHOME_VERBOSE"
	dynamoDBEndpointEnv  = "SMARTHOME_DYNAMODB_ENDPOINT"
	dynamoDBAuthTableEnv = "SMARTHOME_DYNAMODB_AUTH_TABLE"
)

const (
	awsRegionFlag         = "aws.region"
	verboseFlag           = "logging.verbose"
	dynamoDBEndpointFlag  = "aws.dynamodb.endpoint"
	dynamoDBAuthTableFlag = "aws.dynamodb.tables.auth"
)

var (
	c     controller.SmartHomeInterface
	sugar *zap.SugaredLogger
)

type auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

func init() {
	viper.SetDefault(awsRegionFlag, "us-east-1")
	viper.SetDefault(dynamoDBEndpointFlag, "")
	viper.SetDefault(dynamoDBAuthTableFlag, controller.DefaultAuthTable)
	viper.SetDefault(verboseFlag, false)
	viper.BindEnv(awsRegionFlag, awsRegionEnv)
	viper.BindEnv(dynamoDBEndpointFlag, dynamoDBEndpointEnv)
	viper.BindEnv(dynamoDBAuthTableFlag, dynamoDBAuthTableEnv)
	viper.BindEnv(verboseFlag, verboseEnv)

	sugar, err := utils.InitSugaredLogger(viper.GetBool(verboseFlag))

	if err != nil {
		fmt.Printf("error when initializing logger: %s\n", err.Error())
		os.Exit(1)
	}

	region := viper.GetString(awsRegionFlag)
	dynamoDBEndpoint := viper.GetString(dynamoDBEndpointFlag)

	sugar.Infow("creating DynamoDB client", "region", region, "url", dynamoDBEndpoint)
	dynamoClient, err := utils.InitDynamoClient(region, dynamoDBEndpoint)
	if err != nil {
		sugar.Fatalw("error creating DynamoDB client", "error", err.Error())
	}

	c = controller.NewSmartHome(
		controller.SetLogger(sugar),
		controller.SetDynamoDBClient(dynamoClient),
		controller.SetConfig(&controller.SmartHomeConfig{
			AuthTable: viper.GetString(dynamoDBAuthTableFlag),
		}),
	)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	authParams := new(auth)

	if err := json.Unmarshal([]byte(request.Body), &authParams); err != nil {
		return Response{
			Body:       fmt.Sprintf("No valid username or password provided: %s", err.Error()),
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	if err := c.SetCredentials(authParams.Username, authParams.Password); err != nil {
		return Response{
			Body:       fmt.Sprintf("Error saving the credentials in the database: %s", err.Error()),
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("Error saving the credentials in the database: %w", err)
	}

	return Response{
		Body:       "Successfully signed up",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(Handler)
}