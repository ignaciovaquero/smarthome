package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

const envPrefix = "SMARTHOME"

var cfgFile string

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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./smarthome.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("smarthome")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix(envPrefix)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
