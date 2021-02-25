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

const roomParam = "room"

const (
	portEnv                 = "SMARTHOME_SERVER_PORT"
	addressEnv              = "SMARTHOME_LISTEN_ADDRESS"
	jwtSecretEnv            = "SMARTHOME_JWT_SECRET"
	awsRegionEnv            = "SMARTHOME_AWS_REGION"
	verboseEnv              = "SMARTHOME_VERBOSE"
	dynamoDBEndpointEnv     = "SMARTHOME_DYNAMODB_ENDPOINT"
	dynamoDBControlTableEnv = "SMARTHOME_DYNAMODB_CONTROL_PLANE_TABLE"
	dynamoDBOutsideTableEnv = "SMARTHOME_DYNAMODB_TEMPERATURE_OUTSIDE_TABLE"
	dynamoDBInsideTableEnv  = "SMARTHOME_DYNAMODB_TEMPERATURE_INSIDE_TABLE"
)

const (
	portFlag                 = "server.port"
	addressFlag              = "server.address"
	jwtSecretFlag            = "server.jwt.secret"
	awsRegionFlag            = "aws.region"
	verboseFlag              = "logging.verbose"
	dynamoDBEndpointFlag     = "aws.dynamodb.endpoint"
	dynamoDBControlTableFlag = "aws.dynamodb.tables.control"
	dynamoDBOutsideTableFlag = "aws.dynamodb.tables.outside"
	dynamoDBInsideTableFlag  = "aws.dynamodb.tables.inside"
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
	viper.SetDefault(awsRegionFlag, "us-west-1")
	viper.SetDefault(dynamoDBEndpointFlag, "")
	viper.SetDefault(dynamoDBControlTableFlag, controller.DefaultControlPlaneTable)
	viper.SetDefault(dynamoDBOutsideTableFlag, controller.DefaultTempOutsideTable)
	viper.SetDefault(dynamoDBInsideTableFlag, controller.DefaultTempInsideTable)
	viper.SetDefault(verboseFlag, false)
	viper.BindEnv(jwtSecretFlag, jwtSecretEnv)
	viper.BindEnv(awsRegionFlag, awsRegionEnv)
	viper.BindEnv(dynamoDBEndpointFlag, dynamoDBEndpointEnv)
	viper.BindEnv(dynamoDBControlTableFlag, dynamoDBControlTableEnv)
	viper.BindEnv(dynamoDBOutsideTableFlag, dynamoDBOutsideTableEnv)
	viper.BindEnv(dynamoDBInsideTableFlag, dynamoDBInsideTableEnv)
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
			TempOutsideTable:  viper.GetString(dynamoDBOutsideTableFlag),
			TempInsideTable:   viper.GetString(dynamoDBInsideTableFlag),
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

	r := new(controller.RoomOptions)

	if err := json.Unmarshal([]byte(request.Body), &r); err != nil {
		return Response{
			Body:       err.Error(),
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	if r.ThresholdOn >= r.ThresholdOff {
		return Response{
			Body:       "threshold_on should be lower or equal to threshold_off",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	if err := c.SetRoomOptions(room, r); err != nil {
		return Response{
			Body:       fmt.Sprintf("Internal Server Error: %s", err.Error()),
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error setting room options: %w", err)
	}

	return Response{
		Body:       "successfully set room options",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
