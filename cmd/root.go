package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

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
		&configFile, "config", "c", "./config/config.yaml",
		"config file",
	)

}

// InitConfig reads in config file and env vars
func initConfig() {
	config = viper.New()
	config.AddConfigPath("./config/")
	config.SetConfigType("yaml")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	}
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.SetEnvPrefix("pitayaadmin")
	config.AutomaticEnv()

	if err := config.ReadInConfig(); err != nil {
		fmt.Printf("Config file %s failed to load: %s.\n", configFile, err.Error())
		panic("Failed to load config file")
	}

}
