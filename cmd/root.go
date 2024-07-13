package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	okCode   = 0
	failCode = 1
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use: "service",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.BindEnv("storage.host", "STORAGE_HOST")
	viper.BindEnv("storage.port", "STORAGE_PORT")
	viper.BindEnv("storage.dbname", "STORAGE_DBNAME")
	viper.BindEnv("storage.user", "STORAGE_USER")
	viper.BindEnv("storage.password", "STORAGE_PASSWORD")

	viper.BindEnv("notify.host", "NOTIFY_HOST")
	viper.BindEnv("notify.port", "NOTIFY_PORT")
	viper.BindEnv("notify.topic", "NOTIFY_TOPIC")
	viper.BindEnv("notify.partition", "NOTIFY_PARTITION")

	viper.BindEnv("app.host", "APP_HOST")
	viper.BindEnv("app.port", "APP_PORT")

	viper.BindEnv("logger.level", "LOGGER_LEVEL")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
