package cmd

import (
	"fmt"
	"os"

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
)

const (
	portFlag      = "server.port"
	addressFlag   = "server.address"
	jwtSecretFlag = "server.jwt.secret"
)

const apiVersion string = "v1"

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "starts a smarthome server",
	Long: `Starts a smarthome server and listens on the port
	specified in the configuration file, in the command line
	flags or in the corresponding environment variable.`,
	Run: serve,
}

func serve(cmd *cobra.Command, args []string) {
	address := viper.GetString(addressFlag)
	port := viper.GetInt(portFlag)
	jwtSecret := viper.GetString(jwtSecretFlag)
	a := api.NewAPI(api.SetLogger(sugar))

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
	viper.BindPFlag(portFlag, serveCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag(addressFlag, serveCmd.PersistentFlags().Lookup("address"))
	viper.BindEnv(portFlag, portEnv)
	viper.BindEnv(addressFlag, addressEnv)
	viper.BindEnv(jwtSecretFlag, jwtSecretEnv)
}
