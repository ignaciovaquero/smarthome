package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	portFlag    = "server.port"
	addressFlag = "server.address"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "starts a smarthome server",
	Long: `Starts a smarthome server and listens on the port
	specified in the configuration file, in the command line
	flags or in the corresponding environment variable.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")
		fmt.Println("port: ", viper.GetString(portFlag))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().Int("port", 8080, "port where to listen on (default is 8080)")
	serveCmd.PersistentFlags().String("address", "0.0.0.0", "address where to bind to (default is 0.0.0.0)")
	viper.BindPFlag(portFlag, serveCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag(addressFlag, serveCmd.PersistentFlags().Lookup("address"))
	viper.SetEnvPrefix(envPrefix)
	viper.BindEnv(portFlag, "SERVER_PORT")
}
