package cmd

import (
	"flag"

	"github.com/spf13/cobra"
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya-admin/api"
	"github.com/topfreegames/pitaya/serialize/json"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts pitaya admin",
	Long:  `starts pitaya admin api`,
	Run: func(cmd *cobra.Command, args []string) {
		svType := flag.String("type", "admin", "the server type")
		isFrontend := flag.Bool("frontend", false, "if server is frontend")
		flag.Parse()

		go api.Init()

		defer pitaya.Shutdown()

		pitaya.SetSerializer(json.NewSerializer())
		pitaya.Configure(*isFrontend, *svType, pitaya.Cluster, map[string]string{})
		pitaya.Start()
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
}
