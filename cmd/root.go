package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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
