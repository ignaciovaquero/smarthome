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
	headers := map[string]string{}
	if viper.GetString(corsOriginsFlag) != "" {
		headers["Access-Control-Allow-Origin"] = viper.GetString(corsOriginsFlag)
	}

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

	r := new(controller.RoomOptions)

	if err := json.Unmarshal([]byte(request.Body), &r); err != nil {
		return Response{
			Body:       err.Error(),
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
		}, nil
	}

	if r.ThresholdOn > r.ThresholdOff {
		return Response{
			Body:       "threshold_on should be lower or equal to threshold_off",
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
		}, nil
	}

	rooms := []string{room}

	if room == "all" {
		rooms = utils.AllButOne(validRooms, "all")
	}

	for _, roomName := range rooms {
		if err := c.SetRoomOptions(roomName, r); err != nil {
			return Response{
				Body:       fmt.Sprintf("Internal Server Error: %s", err.Error()),
				StatusCode: http.StatusInternalServerError,
				Headers:    headers,
			}, fmt.Errorf("error setting room options for room %s: %w", roomName, err)
		}
	}

	return Response{
		Body:       "successfully set room options",
		StatusCode: http.StatusOK,
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
