package api

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/topfreegames/pitaya-admin/config"
)

// App struct
type App struct {
	Address  string
	Router   *mux.Router
	Server   *http.Server
	Logger   logrus.FieldLogger
	usesGrpc bool
	Config   *config.Config
}

// NewApp constructor
func NewApp(
	host string,
	port int,
	usesGrpc bool,
	config *config.Config,
) (*App, error) {
	a := &App{
		Address:  fmt.Sprintf("%s:%d", host, port),
		Config:   config,
		usesGrpc: usesGrpc,
	}
	a.configureApp()
	return a, nil
}

func (a *App) configureApp() {
	a.Router = a.getRouter()
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
