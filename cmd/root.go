package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ConfigFile is the configuration file used for running a command
var ConfigFile string

// Config is the app configuration
var config *viper.Viper

// RootCmd is the root command
var RootCmd = &cobra.Command{
	Use:   "pitaya admin",
	Short: "Pitaya admin is an admin for pitaya's servers",
	Long:  "Use pitaya admin to manage your servers",
}

// Execute runs RootCmd to initialize pitaya admin
func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVarP(
		&ConfigFile, "config", "c", "./config/config.yaml",
		"config file",
	)

}

// InitConfig reads in config file and env vars
func initConfig() {
	config = viper.New()
	config.AddConfigPath("./config/")
	config.SetConfigType("yaml")
	if ConfigFile != "" {
		viper.SetConfigFile(ConfigFile)
	}
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.SetEnvPrefix("pitayaadmin")
	config.AutomaticEnv()

	if err := config.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %s", config.ConfigFileUsed())
	}

}
