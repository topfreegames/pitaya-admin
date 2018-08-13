package main

import (
	"flag"

	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya-admin/api"
	"github.com/topfreegames/pitaya/serialize/json"
)

func main() {
	svType := flag.String("type", "specAdmin", "the server type")
	isFrontend := flag.Bool("frontend", false, "if server is frontend")
	flag.Parse()

	go api.Init()

	defer pitaya.Shutdown()

	pitaya.SetSerializer(json.NewSerializer())
	pitaya.Configure(*isFrontend, *svType, pitaya.Cluster, map[string]string{})
	pitaya.Start()
}
