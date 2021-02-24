package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/igvaquero18/smarthome/api"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	portEnv                 = "SMARTHOME_SERVER_PORT"
	addressEnv              = "SMARTHOME_LISTEN_ADDRESS"
	jwtSecretEnv            = "SMARTHOME_JWT_SECRET"
	awsRegionEnv            = "SMARTHOME_AWS_REGION"
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
	dynamoDBEndpointFlag     = "aws.dynamodb.endpoint"
	dynamoDBControlTableFlag = "aws.dynamodb.tables.control"
	dynamoDBOutsideTableFlag = "aws.dynamodb.tables.outside"
	dynamoDBInsideTableFlag  = "aws.dynamodb.tables.inside"
)

const apiVersion string = "v1"

// serveCmd represents the serve command
var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "starts a smarthome server",
		Long: `Starts a smarthome server and listens on the port
	specified in the configuration file, in the command line
	flags or in the corresponding environment variable.`,
		Run: serve,
	}
)

func serve(cmd *cobra.Command, args []string) {
	region := viper.GetString(awsRegionFlag)
	jwtSecret := viper.GetString(jwtSecretFlag)
	address := viper.GetString(addressFlag)
	port := viper.GetInt(portFlag)
	dynamoDBEndpoint := viper.GetString(dynamoDBEndpointFlag)
	dynamoDBControlTable := viper.GetString(dynamoDBControlTableFlag)

	sugar.Infow("creating DynamoDB client", "region", region, "url", dynamoDBEndpoint)
	dynamoClient, err := initDynamoClient(region, dynamoDBEndpoint)
	if err != nil {
		sugar.Fatalw("error creating DynamoDB client", "error", err.Error())
	}

	s := api.NewSmartHome(
		api.SetLogger(sugar),
		api.SetDynamoDBClient(dynamoClient),
		api.SetConfig(&api.SmartHomeConfig{
			ControlPlaneTable: dynamoDBControlTable,
		}),
	)

	sugar.Infow("starting server", "address", address, "port", port)

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"level":"info","ts":"${time_unix}","id":"${id}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status},"error":"${error}","latency":${latency},` +
			`"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		Output: os.Stdout,
	}))

	if jwtSecret != "" {
		e.Use(middleware.JWT([]byte(jwtSecret)))
	} else {
		sugar.Warn("no jwt secret provided, disabling authentication")
	}

	room := e.Group(fmt.Sprintf("%s/room", apiVersion))
	room.POST("/:room", s.SetRoomOptions)
	room.GET("/:room", s.GetRoomOptions)
	p := prometheus.NewPrometheus("smarthome", nil)
	p.Use(e)

	sugar.Infow("starting server", "address", address, "port", port)
	sugar.Fatal(e.Start(fmt.Sprintf("%s:%d", address, port)))
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntP("port", "p", 8080, "port where to listen on")
	serveCmd.Flags().StringP("address", "a", "0.0.0.0", "address where to bind to")
	serveCmd.Flags().StringP("aws-region", "r", "us-west-1", "AWS region for DynamoDB")
	serveCmd.Flags().StringP("dynamodb-endpoint", "d", "", "DynamoDB endpoint")
	serveCmd.Flags().String("dynamodb-control-table", "ControlPlane", "DynamoDB Control Plane table name")
	serveCmd.Flags().String("dynamodb-outside-table", "TemperatureOutside", "DynamoDB Temperature Outside table name")
	serveCmd.Flags().String("dynamodb-inside-table", "TemperatureInside", "DynamoDB Temperature Inside table name")
	viper.BindPFlag(portFlag, serveCmd.Flags().Lookup("port"))
	viper.BindPFlag(addressFlag, serveCmd.Flags().Lookup("address"))
	viper.BindPFlag(awsRegionFlag, serveCmd.Flags().Lookup("aws-region"))
	viper.BindPFlag(dynamoDBEndpointFlag, serveCmd.Flags().Lookup("dynamodb-endpoint"))
	viper.BindPFlag(dynamoDBControlTableFlag, serveCmd.Flags().Lookup("dynamodb-control-table"))
	viper.BindPFlag(dynamoDBOutsideTableFlag, serveCmd.Flags().Lookup("dynamodb-outside-table"))
	viper.BindPFlag(dynamoDBInsideTableFlag, serveCmd.Flags().Lookup("dynamodb-inside-table"))
	viper.BindEnv(portFlag, portEnv)
	viper.BindEnv(addressFlag, addressEnv)
	viper.BindEnv(jwtSecretFlag, jwtSecretEnv)
	viper.BindEnv(awsRegionFlag, awsRegionEnv)
	viper.BindEnv(dynamoDBEndpointFlag, dynamoDBEndpointEnv)
	viper.BindEnv(dynamoDBControlTableFlag, dynamoDBControlTableEnv)
	viper.BindEnv(dynamoDBOutsideTableFlag, dynamoDBOutsideTableEnv)
	viper.BindEnv(dynamoDBInsideTableFlag, dynamoDBInsideTableEnv)
}

func initDynamoClient(region, url string) (*dynamodb.Client, error) {
	var cfg aws.Config
	var err error

	if url == "" {
		cfg, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(region),
		)
	} else {
		customResolver := aws.EndpointResolverFunc(func(service, awsRegion string) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           url,
				SigningRegion: awsRegion,
			}, nil
		})
		cfg, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(region),
			config.WithEndpointResolver(customResolver),
		)
	}
	if err != nil {
		return nil, fmt.Errorf("error loading aws configuration: %w", err)
	}
	return dynamodb.NewFromConfig(cfg), nil
}
