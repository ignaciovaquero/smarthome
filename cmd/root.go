package cmd

import (
	"fmt"
	"os"

	"code.hq.twilio.com/twilio/owl-multiaccount-server/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/spf13/viper"
)

const (
	configFlag  = "config"
	verboseFlag = "logging.verbose"
)

var (
	sugar *zap.SugaredLogger
)

var rootCmd = &cobra.Command{
	Use:   "smarthome",
	Short: "an application for centrally controlling a Homekit based smart home",
	Long: `SmartHome is an API that controlls my Homekit based Smart Home.

	It allows to set whether we want to manually control the temperature
	of the home, or whether we want it to be automatically set based on
	the smart thermometers.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("config", "c", "", `config file (default "./smarthome.yaml")`)
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose logging (default false)")
	viper.BindPFlag(configFlag, rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag(verboseFlag, rootCmd.PersistentFlags().Lookup("verbose"))

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	sugar = logger.Sugar()
}

func initConfig() {
	var err error

	if cfgFile := viper.GetString(configFlag); cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("smarthome")
		viper.SetConfigType("yml")
	}
	if err = viper.ReadInConfig(); err == nil {
		fmt.Printf("using config file %s\n", viper.ConfigFileUsed())
	}

	if sugar, err = utils.InitSugaredLogger(viper.GetBool(verboseFlag)); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
