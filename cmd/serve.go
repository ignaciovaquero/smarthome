package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/igvaquero18/smarthome/api"
	"github.com/igvaquero18/smarthome/controller"
	"github.com/igvaquero18/smarthome/utils"
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
	jwtExpirationEnv        = "SMARTHOME_JWT_EXPIRATION"
	awsRegionEnv            = "SMARTHOME_AWS_REGION"
	corsOriginsEnv          = "SMARTHOME_CORS_ORIGINS"
	dynamoDBEndpointEnv     = "SMARTHOME_DYNAMODB_ENDPOINT"
	dynamoDBAuthTableEnv    = "SMARTHOME_DYNAMODB_AUTH_TABLE"
	dynamoDBControlTableEnv = "SMARTHOME_DYNAMODB_CONTROL_PLANE_TABLE"
	dynamoDBOutsideTableEnv = "SMARTHOME_DYNAMODB_TEMPERATURE_OUTSIDE_TABLE"
	dynamoDBInsideTableEnv  = "SMARTHOME_DYNAMODB_TEMPERATURE_INSIDE_TABLE"
)

const (
	portFlag                 = "server.port"
	addressFlag              = "server.address"
	jwtSecretFlag            = "server.jwt.secret"
	jwtExpirationFlag        = "server.jwt.expiration"
	awsRegionFlag            = "aws.region"
	corsOriginsFlag          = "cors.origins"
	dynamoDBEndpointFlag     = "aws.dynamodb.endpoint"
	dynamoDBAuthTableFlag    = "aws.dynamodb.tables.auth"
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
	origins := strings.Split(viper.GetString(corsOriginsFlag), " ")
	dynamoDBEndpoint := viper.GetString(dynamoDBEndpointFlag)
	dynamoDBAuthTable := viper.GetString(dynamoDBAuthTableFlag)
	dynamoDBControlTable := viper.GetString(dynamoDBControlTableFlag)
	dynamoDBOutsiteTable := viper.GetString(dynamoDBOutsideTableFlag)
	dynamoDBInsiteTable := viper.GetString(dynamoDBInsideTableFlag)

	sugar.Infow("creating DynamoDB client", "region", region, "url", dynamoDBEndpoint)
	dynamoClient, err := utils.InitDynamoClient(region, dynamoDBEndpoint)
	if err != nil {
		sugar.Fatalw("error creating DynamoDB client", "error", err.Error())
	}

	expiration, err := time.ParseDuration(viper.GetString(jwtExpirationFlag))

	if err != nil {
		sugar.Fatalw("invalid parameters for the JWT expiration time", "expiration", viper.GetString(jwtExpirationFlag))
	}

	s := api.NewClient(
		api.JWTConfig{
			JWTSecret:     jwtSecret,
			JWTExpiration: expiration,
		},
		controller.NewSmartHome(
			controller.SetLogger(sugar),
			controller.SetDynamoDBClient(dynamoClient),
			controller.SetConfig(&controller.SmartHomeConfig{
				AuthTable:         dynamoDBAuthTable,
				ControlPlaneTable: dynamoDBControlTable,
				TempOutsideTable:  dynamoDBOutsiteTable,
				TempInsideTable:   dynamoDBInsiteTable,
			}),
		),
	)

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"ts":"${time_unix}","id":"${id}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status},"error":"${error}","latency":${latency},` +
			`"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		Output: os.Stdout,
	}))

	if len(origins) > 0 {
		if err := utils.ValidateURLsFromArray(origins); err != nil {
			sugar.Fatalw("invalid CORS URLs provided", "error", err.Error())
		}
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: origins,
		}))
	}

	room := e.Group(fmt.Sprintf("%s/room", apiVersion))
	if jwtSecret != "" {
		room.Use(middleware.JWT([]byte(jwtSecret)))
		e.POST(fmt.Sprintf("%s/login", apiVersion), s.Login)
		e.POST(fmt.Sprintf("%s/signup", apiVersion), s.SignUp)
		e.DELETE(fmt.Sprintf("%s/user", apiVersion), s.DeleteUser)
	} else {
		sugar.Warn("no jwt secret provided, disabling authentication")
	}
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
	serveCmd.Flags().StringP("aws-region", "r", "us-east-1", "AWS region for DynamoDB")
	serveCmd.Flags().StringP("dynamodb-endpoint", "d", "", "DynamoDB endpoint")
	serveCmd.Flags().String("dynamodb-auth-table", controller.DefaultAuthTable, "DynamoDB Authentication table name")
	serveCmd.Flags().String("dynamodb-control-table", controller.DefaultControlPlaneTable, "DynamoDB Control Plane table name")
	serveCmd.Flags().String("dynamodb-outside-table", controller.DefaultTempOutsideTable, "DynamoDB Temperature Outside table name")
	serveCmd.Flags().String("dynamodb-inside-table", controller.DefaultTempInsideTable, "DynamoDB Temperature Inside table name")
	serveCmd.Flags().String("jwt-expiration", "1h", "Expiration of JWT token. See https://golang.org/pkg/time/#ParseDuration for an example of how to set this parameter")
	serveCmd.Flags().String("cors-origins", "", "Space-separated list of CORS Origin URLs")
	viper.BindPFlag(portFlag, serveCmd.Flags().Lookup("port"))
	viper.BindPFlag(addressFlag, serveCmd.Flags().Lookup("address"))
	viper.BindPFlag(awsRegionFlag, serveCmd.Flags().Lookup("aws-region"))
	viper.BindPFlag(dynamoDBEndpointFlag, serveCmd.Flags().Lookup("dynamodb-endpoint"))
	viper.BindPFlag(dynamoDBAuthTableFlag, serveCmd.Flags().Lookup("dynamodb-auth-table"))
	viper.BindPFlag(dynamoDBControlTableFlag, serveCmd.Flags().Lookup("dynamodb-control-table"))
	viper.BindPFlag(dynamoDBOutsideTableFlag, serveCmd.Flags().Lookup("dynamodb-outside-table"))
	viper.BindPFlag(dynamoDBInsideTableFlag, serveCmd.Flags().Lookup("dynamodb-inside-table"))
	viper.BindPFlag(jwtExpirationFlag, serveCmd.Flags().Lookup("jwt-expiration"))
	viper.BindPFlag(corsOriginsFlag, serveCmd.Flags().Lookup("cors-origins"))
	viper.BindEnv(portFlag, portEnv)
	viper.BindEnv(addressFlag, addressEnv)
	viper.BindEnv(jwtSecretFlag, jwtSecretEnv)
	viper.BindEnv(jwtExpirationFlag, jwtExpirationEnv)
	viper.BindEnv(awsRegionFlag, awsRegionEnv)
	viper.BindEnv(corsOriginsFlag, corsOriginsEnv)
	viper.BindEnv(dynamoDBEndpointFlag, dynamoDBEndpointEnv)
	viper.BindEnv(dynamoDBAuthTableFlag, dynamoDBAuthTableEnv)
	viper.BindEnv(dynamoDBControlTableFlag, dynamoDBControlTableEnv)
	viper.BindEnv(dynamoDBOutsideTableFlag, dynamoDBOutsideTableEnv)
	viper.BindEnv(dynamoDBInsideTableFlag, dynamoDBInsideTableEnv)
}
