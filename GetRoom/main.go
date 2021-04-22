package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/igvaquero18/smarthome/controller"
	"github.com/igvaquero18/smarthome/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	awsRegionEnv            = "SMARTHOME_AWS_REGION"
	verboseEnv              = "SMARTHOME_VERBOSE"
	dynamoDBEndpointEnv     = "SMARTHOME_DYNAMODB_ENDPOINT"
	dynamoDBControlTableEnv = "SMARTHOME_DYNAMODB_CONTROL_PLANE_TABLE"
)

const (
	awsRegionFlag            = "aws.region"
	verboseFlag              = "logging.verbose"
	dynamoDBEndpointFlag     = "aws.dynamodb.endpoint"
	dynamoDBControlTableFlag = "aws.dynamodb.tables.control"
)

var (
	validRooms = []string{"all", "bedroom", "livingroom"}
	c          controller.SmartHomeInterface
	sugar      *zap.SugaredLogger
)

type validRoom string

func (r validRoom) isValid() bool {
	for _, room := range validRooms {
		if string(r) == room {
			return true
		}
	}
	return false
}

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

func init() {
	viper.SetDefault(awsRegionFlag, "us-east-1")
	viper.SetDefault(dynamoDBEndpointFlag, "")
	viper.SetDefault(dynamoDBControlTableFlag, controller.DefaultControlPlaneTable)
	viper.SetDefault(verboseFlag, false)
	viper.BindEnv(awsRegionFlag, awsRegionEnv)
	viper.BindEnv(dynamoDBEndpointFlag, dynamoDBEndpointEnv)
	viper.BindEnv(dynamoDBControlTableFlag, dynamoDBControlTableEnv)
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
			ControlPlaneTable: viper.GetString(dynamoDBControlTableFlag),
		}),
	)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	room := request.PathParameters["room"]

	if !validRoom(room).isValid() {
		return Response{
			Body:       "Invalid room name",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	if room == "all" {
		rooms := utils.AllButOne(validRooms, "all")
		roomOpts := []map[string]types.AttributeValue{}
		for _, roomName := range rooms {
			item, err := c.GetRoomOptions(roomName)
			if err != nil {
				return Response{
					Body:       fmt.Sprintf("Internal Server Error: %s", err.Error()),
					StatusCode: http.StatusInternalServerError,
				}, fmt.Errorf("error getting item from DynamoDB: %w", err)
			}
			if item == nil {
				continue
			}
			roomOpts = append(roomOpts, item)
		}

		if len(roomOpts) == 0 {
			return Response{
				Body:       "Not found",
				StatusCode: http.StatusNotFound,
			}, nil
		}
		body, err := json.Marshal(roomOpts)
		if err != nil {
			return Response{
				Body: fmt.Sprintf("Internal Server Error: %s", err.Error()),
			}, fmt.Errorf("error marshalling response: %w", err)
		}
		return Response{
			Body:       string(body),
			StatusCode: http.StatusOK,
		}, nil
	}

	item, err := c.GetRoomOptions(room)
	if err != nil {
		return Response{
			Body:       fmt.Sprintf("Internal Server Error: %s", err.Error()),
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error getting item from DynamoDB: %w", err)
	}

	if item == nil {
		return Response{
			Body:       "Not found",
			StatusCode: http.StatusNotFound,
		}, nil
	}

	body, err := json.Marshal(item)
	if err != nil {
		return Response{
			Body:       fmt.Sprintf("Internal Server Error: %s", err.Error()),
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error marshalling response: %w", err)
	}

	return Response{
		Body:       string(body),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
