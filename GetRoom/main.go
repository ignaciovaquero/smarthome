package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/igvaquero18/smarthome/controller"
	"github.com/igvaquero18/smarthome/utils"
	"github.com/spf13/viper"
)

const (
	jwtSecretEnv            = "SMARTHOME_JWT_SECRET"
	awsRegionEnv            = "SMARTHOME_AWS_REGION"
	verboseEnv              = "SMARTHOME_VERBOSE"
	corsOriginsEnv          = "SMARTHOME_CORS_ORIGINS"
	dynamoDBEndpointEnv     = "SMARTHOME_DYNAMODB_ENDPOINT"
	dynamoDBControlTableEnv = "SMARTHOME_DYNAMODB_CONTROL_PLANE_TABLE"
)

const (
	jwtSecretFlag            = "server.jwt.secret"
	awsRegionFlag            = "aws.region"
	verboseFlag              = "logging.verbose"
	corsOriginsFlag          = "cors.origins"
	dynamoDBEndpointFlag     = "aws.dynamodb.endpoint"
	dynamoDBControlTableFlag = "aws.dynamodb.tables.control"
)

var (
	validRooms = []string{"all", "bedroom", "livingroom"}
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

// RoomOptions is a struct that represents the options available for a room
type RoomOptions struct {
	Name         string  `json:"name,omitempty"`
	Enabled      bool    `json:"enabled"`
	ThresholdOn  float32 `json:"threshold_on"`
	ThresholdOff float32 `json:"threshold_off"`
}

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

func init() {
	viper.SetDefault(jwtSecretFlag, "")
	viper.SetDefault(awsRegionFlag, "us-east-3")
	viper.SetDefault(verboseFlag, false)
	viper.SetDefault(corsOriginsFlag, "")
	viper.SetDefault(dynamoDBEndpointFlag, "")
	viper.SetDefault(dynamoDBControlTableFlag, controller.DefaultControlPlaneTable)
	viper.BindEnv(jwtSecretFlag, jwtSecretEnv)
	viper.BindEnv(awsRegionFlag, awsRegionEnv)
	viper.BindEnv(verboseFlag, verboseEnv)
	viper.BindEnv(corsOriginsFlag, corsOriginsEnv)
	viper.BindEnv(dynamoDBEndpointFlag, dynamoDBEndpointEnv)
	viper.BindEnv(dynamoDBControlTableFlag, dynamoDBControlTableEnv)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	headers := map[string]string{}
	if viper.GetString(corsOriginsFlag) != "" {
		headers["Access-Control-Allow-Origin"] = viper.GetString(corsOriginsFlag)
	}

	sugar, err := utils.InitSugaredLogger(viper.GetBool(verboseFlag))

	if err != nil {
		return Response{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Internal Server Error: %s", err.Error()),
			Headers:    headers,
		}, fmt.Errorf("Error when initializing logger: %w", err)
	}

	region := viper.GetString(awsRegionFlag)
	dynamoDBEndpoint := viper.GetString(dynamoDBEndpointFlag)

	sugar.Infow("creating DynamoDB client", "region", region, "url", dynamoDBEndpoint)
	dynamoClient, err := utils.InitDynamoClient(region, dynamoDBEndpoint)
	if err != nil {
		return Response{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Internal Server Error: %s", err.Error()),
			Headers:    headers,
		}, fmt.Errorf("Error when creating DynamoDB client: %w", err)
	}

	var c controller.SmartHomeInterface = controller.NewSmartHome(
		controller.SetLogger(sugar),
		controller.SetDynamoDBClient(dynamoClient),
		controller.SetConfig(&controller.SmartHomeConfig{
			ControlPlaneTable: viper.GetString(dynamoDBControlTableFlag),
		}),
	)

	if err := utils.ValidateTokenFromHeader(request.Headers["Authorization"], viper.GetString(jwtSecretFlag)); err != nil {
		return Response{
			Body:       fmt.Sprintf("Authentication failure: %s", err.Error()),
			StatusCode: http.StatusForbidden,
			Headers:    headers,
		}, nil
	}

	room := request.PathParameters["room"]

	if !validRoom(room).isValid() {
		return Response{
			Body:       "Invalid room name",
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
		}, nil
	}

	if room == "all" {
		rooms := utils.AllButOne(validRooms, "all")
		roomOpts := []RoomOptions{}
		for _, roomName := range rooms {
			item, err := c.GetRoomOptions(roomName)
			if err != nil {
				return Response{
					Body:       fmt.Sprintf("Internal Server Error: %s", err.Error()),
					StatusCode: http.StatusInternalServerError,
					Headers:    headers,
				}, fmt.Errorf("Error getting item from DynamoDB: %w", err)
			}
			if item == nil {
				continue
			}
			roomOpt := RoomOptions{Name: roomName}
			if err = attributevalue.UnmarshalMap(item, &roomOpt); err != nil {
				return Response{
					StatusCode: http.StatusInternalServerError,
					Body:       fmt.Sprintf("Internal Server error: %s", err.Error()),
					Headers:    headers,
				}, fmt.Errorf("Error unmarshalling DynamoDB item: %w", err)
			}
			roomOpts = append(roomOpts, roomOpt)
		}

		if len(roomOpts) == 0 {
			return Response{
				Body:       "No rooms were found",
				StatusCode: http.StatusNotFound,
				Headers:    headers,
			}, nil
		}
		body, err := json.Marshal(roomOpts)
		if err != nil {
			return Response{
				StatusCode: http.StatusInternalServerError,
				Body:       fmt.Sprintf("Internal Server Error: %s", err.Error()),
				Headers:    headers,
			}, fmt.Errorf("error marshalling response: %w", err)
		}
		return Response{
			Body:       string(body),
			StatusCode: http.StatusOK,
			Headers:    headers,
		}, nil
	}

	item, err := c.GetRoomOptions(room)
	if err != nil {
		return Response{
			Body:       fmt.Sprintf("Internal Server Error: %s", err.Error()),
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
		}, fmt.Errorf("error getting item from DynamoDB: %w", err)
	}

	if item == nil {
		return Response{
			Body:       fmt.Sprintf("Room %s not found", room),
			StatusCode: http.StatusNotFound,
			Headers:    headers,
		}, nil
	}

	roomOpt := RoomOptions{Name: room}
	if err = attributevalue.UnmarshalMap(item, &roomOpt); err != nil {
		return Response{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Internal Server Error: %s", err.Error()),
			Headers:    headers,
		}, fmt.Errorf("Error unmarshalling DynamoDB item: %w", err)
	}

	body, err := json.Marshal(roomOpt)
	if err != nil {
		return Response{
			Body:       fmt.Sprintf("Internal Server Error: %s", err.Error()),
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
		}, fmt.Errorf("error marshalling response: %w", err)
	}

	return Response{
		Body:       string(body),
		StatusCode: http.StatusOK,
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
