package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya-admin/api"
	"github.com/topfreegames/pitaya/cluster"
	"github.com/topfreegames/pitaya/modules"
	"github.com/topfreegames/pitaya/serialize/json"
	"github.com/topfreegames/pitaya/serialize/protobuf"
)

var bind string
var isFrontend bool
var port int
var rpcServerPort string
var svType string
var usesGrpc bool
var usesJSON bool

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

		app, err := api.NewApp(bind, port, config)

		if err != nil {
			cmdL.Fatal(err)
		}

		defer pitaya.Shutdown()
		app.Init()

		if usesJSON {
			pitaya.SetSerializer(json.NewSerializer())
		} else {
			pitaya.SetSerializer(protobuf.NewSerializer())
		}

		if usesGrpc {
			if rpcServerPort == "" {
				cmdL.Fatal("cant use grpc without a port, specify with -r")
			}

			meta := map[string]string{
				"grpc-host": "0.0.0.0",
				"grpc-port": rpcServerPort,
			}

			config.Set("pitaya.cluster.rpc.server.grpc.port", rpcServerPort)

			pitaya.Configure(isFrontend, svType, pitaya.Cluster, meta, config)

			bs := modules.NewETCDBindingStorage(pitaya.GetServer(), pitaya.GetConfig())
			pitaya.RegisterModule(bs, "bindingsStorage")

			gs, err := cluster.NewGRPCServer(pitaya.GetConfig(), pitaya.GetServer(), pitaya.GetMetricsReporters())
			if err != nil {
				cmdL.Fatal(err)
			}

			gc, err := cluster.NewGRPCClient(pitaya.GetConfig(), pitaya.GetServer(), pitaya.GetMetricsReporters(), bs)
			if err != nil {
				cmdL.Fatal(err)
			}
			pitaya.SetRPCServer(gs)
			pitaya.SetRPCClient(gc)

		} else {
			pitaya.Configure(isFrontend, svType, pitaya.Cluster, map[string]string{}, config)
		}

		pitaya.Start()
	},
}

func init() {
	startCmd.Flags().BoolVarP(&usesJSON, "usesJSON", "j", true, "if server uses json or not")
	startCmd.Flags().BoolVarP(&usesGrpc, "usesGrpc", "g", false, "if server uses or not grpc")
	startCmd.Flags().StringVarP(&rpcServerPort, "rpcServerPort", "r", "", "the port that grpc server will listen")
	startCmd.Flags().BoolVar(&isFrontend, "isFrontend", false, "if server is frontend")
	startCmd.Flags().StringVar(&svType, "type", "admin", "the server type")
	startCmd.Flags().StringVarP(&bind, "bind", "b", "0.0.0.0", "bind address")
	startCmd.Flags().IntVarP(&port, "port", "p", 8000, "bind port")
	RootCmd.AddCommand(startCmd)
}
