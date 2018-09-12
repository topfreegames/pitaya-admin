package testing

import (
	"time"

	"github.com/spf13/viper"
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya-admin/api"
	"github.com/topfreegames/pitaya/serialize/json"
)

// TestApp is the default testing app
var TestApp *api.App

// IsConf specifies if app is configured
var IsConf bool

// ConfApp configure a default test app
func ConfApp() {

	conf := viper.New()
	TestApp, _ = api.NewApp("0.0.0.0", 8000, conf)

	TestApp.Init()
	pitaya.GetConfig()
	pitaya.SetSerializer(json.NewSerializer())
	pitaya.Configure(false, "pitaya-admin", pitaya.Cluster, map[string]string{})
	go pitaya.Start()
	time.Sleep(500 * time.Millisecond)
	IsConf = true
}
