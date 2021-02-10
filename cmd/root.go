package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/spf13/viper"
)

const (
	configFlag  = "config"
	verboseFlag = "verbose"
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
	rootCmd.PersistentFlags().StringP(configFlag, "c", "", `config file (default "./smarthome.yaml")`)
	rootCmd.PersistentFlags().BoolP(verboseFlag, "v", false, "verbose logging (default false)")
	viper.BindPFlag(configFlag, rootCmd.PersistentFlags().Lookup(configFlag))
	viper.BindPFlag(verboseFlag, rootCmd.PersistentFlags().Lookup(verboseFlag))
	if err := initLogger(); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func initLogger() error {
	var zl *zap.Logger
	cfg := zap.Config{
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	if viper.GetBool(verboseFlag) {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	zl, err := cfg.Build()

	if err != nil {
		return fmt.Errorf("error when initializing logger: %w", err)
	}

	sugar = zl.Sugar()
	sugar.Debug("logger initialization successful")
	return nil
}

func initConfig() {
	if cfgFile := viper.GetString(configFlag); cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("smarthome")
		viper.SetConfigType("yml")
	}
	if err := viper.ReadInConfig(); err == nil {
		sugar.Debugw("using config file", "config_file", viper.ConfigFileUsed())
	}
}
