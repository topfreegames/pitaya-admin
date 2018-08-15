package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya-admin/api"
	"github.com/topfreegames/pitaya/serialize/json"
)

var isFrontend bool
var svType string
var bind string
var port int

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts pitaya admin",
	Long:  "starts pitaya admin api",
	Run: func(cmd *cobra.Command, args []string) {

		var log = logrus.New()
		cmdL := log.WithFields(logrus.Fields{
			"source":    "startCmd",
			"operation": "Run",
			"bind":      bind,
			"port":      port,
		})

		cmdL.Info("starting pitaya admin")
		app, err := api.NewApp(bind, port)

		if err != nil {
			cmdL.Fatal(err)
		}

		defer pitaya.Shutdown()
		app.Init()

		pitaya.SetSerializer(json.NewSerializer())
		pitaya.Configure(isFrontend, svType, pitaya.Cluster, map[string]string{})
		pitaya.Start()
	},
}

func init() {
	startCmd.Flags().BoolVar(&isFrontend, "isFrontend", false, "if server is frontend")
	startCmd.Flags().StringVar(&svType, "type", "admin", "the server type")
	startCmd.Flags().StringVarP(&bind, "bind", "b", "0.0.0.0", "bind address")
	startCmd.Flags().IntVarP(&port, "port", "p", 8000, "bind port")
	RootCmd.AddCommand(startCmd)
}
