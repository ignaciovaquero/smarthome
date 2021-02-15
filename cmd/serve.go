package cmd

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

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
	portEnv      = "SMARTHOME_SERVER_PORT"
	addressEnv   = "SMARTHOME_LISTEN_ADDRESS"
	jwtSecretEnv = "SMARTHOME_JWT_SECRET"
	awsRegionEnv = "SMARTHOME_AWS_REGION"
)

const (
	portFlag      = "server.port"
	addressFlag   = "server.address"
	jwtSecretFlag = "server.jwt.secret"
	awsRegionFlag = "aws.region"
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
	address := viper.GetString(addressFlag)
	port := viper.GetInt(portFlag)
	jwtSecret := viper.GetString(jwtSecretFlag)
	region := viper.GetString(awsRegionFlag)
	a := api.NewAPI(api.SetLogger(sugar))
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolver(aws.NewConfig().EndpointResolver),
	)
	if err != nil {
		sugar.Fatalw("error loading aws configuration", "error", err.Error())
	}

	sugar.Debugw("creating DynamoDB client", "region", region)
	dynamoClient := dynamodb.NewFromConfig(cfg)
	waitTime := 7 * time.Minute
	dynamoTables := []string{"Home", "Indoor Temperature", "Outdoor Temperature"}

	errCh := make(chan error)
	wgDone := make(chan bool, 1)
	wg := &sync.WaitGroup{}
	wg.Add(len(dynamoTables))

	sugar.Infow("creating DynamoDB tables...", "num_tables", len(dynamoTables))
	for _, dynamoTable := range dynamoTables {
		go func(table string) {
			sugar.Debugw("creating table", "table", table)
			e := createTable(dynamoClient, table, waitTime)
			if e != nil {
				errCh <- fmt.Errorf(`error creating DynamoDB table "%s": %w`, table, err)
			}
			wg.Done()
		}(dynamoTable)
	}

	go func() {
		wg.Wait()
		close(wgDone)
	}()

	select {
	case err = <-errCh:
		sugar.Fatalw("error creating tables", "error", err.Error())
	case <-wgDone:
		break
	}

	sugar.Info("successfully created DynamoDb tables")
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

	autotune := e.Group(fmt.Sprintf("%s/autoadjust", apiVersion))
	autotune.POST("/:room", a.AutoAdjustTemperature)
	p := prometheus.NewPrometheus("smarthome", nil)
	p.Use(e)

	sugar.Fatal(e.Start(fmt.Sprintf("%s:%d", address, port)))
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().IntP("port", "p", 8080, "port where to listen on")
	serveCmd.PersistentFlags().StringP("address", "a", "0.0.0.0", "address where to bind to")
	serveCmd.PersistentFlags().StringP("aws-region", "r", "us-west-1", "AWS region for DynamoDB")
	viper.BindPFlag(portFlag, serveCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag(addressFlag, serveCmd.PersistentFlags().Lookup("address"))
	viper.BindPFlag(awsRegionFlag, serveCmd.PersistentFlags().Lookup("aws-region"))
	viper.BindEnv(portFlag, portEnv)
	viper.BindEnv(addressFlag, addressEnv)
	viper.BindEnv(jwtSecretFlag, jwtSecretEnv)
	viper.BindEnv(awsRegionFlag, awsRegionEnv)
}

func initDynamoClient() {

}

func createTable(client *dynamodb.Client, tableName string, waitTime time.Duration) error {
	waiter := dynamodb.NewTableExistsWaiter(client)
	params := &dynamodb.DescribeTableInput{
		TableName: &tableName,
	}
	return waiter.Wait(context.TODO(), params, waitTime)
}
