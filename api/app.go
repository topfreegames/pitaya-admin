package api

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
)

// App struct
type App struct {
	Address string
	Router  *mux.Router
	Server  *http.Server
	Logger  logrus.FieldLogger
	Config  *viper.Viper
}

// NewApp constructor
func NewApp(
	host string,
	port int,
	config *viper.Viper,
) (*App, error) {
	a := &App{
		Address: fmt.Sprintf("%s:%d", host, port),
		Config:  config,
	}
	a.configureApp()
	return a, nil
}

func (a *App) configureApp() {
	a.Router = a.getRouter()
	a.Config.SetDefault("routes.docs", "connectorremote.docs")
	a.Config.SetDefault("routes.protos", "connectorremote.proto")
	a.Config.SetDefault("request.readdeadline", "15s")
	a.Config.SetDefault("request.whitelist", []string{"0.0.0.0:8000"})
	a.configureServer()
}

func (a *App) configureServer() {
	a.Server = &http.Server{
		Addr:    a.Address,
		Handler: a.Router,
	}
}

func (a *App) getRouter() *mux.Router {
	router := mux.NewRouter()
	router.Handle("/docs", NewDocsHandler(a)).Methods("GET")
	router.Handle("/request", NewRequestHandler(a)).Methods("GET")
	router.Handle("/rpc", NewRPCHandler(a)).Methods("POST")
	router.Handle("/servers", NewServersHandler(a)).Methods("GET")
	router.Handle("/user/kick", NewKickUserHandler(a)).Methods("POST")
	router.Handle("/user/push", NewPushToUsersHandler(a)).Methods("POST")
	return router
}

// Init starts the app
func (a *App) Init() {
	go a.Server.ListenAndServe()
}

// GetRemoteDocsRoute gets the route for the autodoc handler
func (a *App) GetRemoteDocsRoute() string {
	return a.Config.GetString("routes.docs")
}

// GetRemoteProtosRoute gets the route for the protos handler
func (a *App) GetRemoteProtosRoute() string {
	return a.Config.GetString("routes.protos")
}
