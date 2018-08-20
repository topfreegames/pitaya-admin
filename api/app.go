package api

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

// App struct
type App struct {
	Address  string
	Router   *mux.Router
	Server   *http.Server
	Logger   logrus.FieldLogger
	usesGrpc bool
}

// NewApp constructor
func NewApp(
	host string,
	port int,
	usesGrpc bool,
) (*App, error) {
	a := &App{
		Address:  fmt.Sprintf("%s:%d", host, port),
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
	router.Handle("/servers", NewServersHandler(a)).Methods("GET")
	router.Handle("/user/push", NewPushToUsersHandler(a)).Methods("POST")
	router.Handle("/user/kick", NewKickUserHandler(a)).Methods("POST")
	return router
}

// Init starts the app
func (a *App) Init() {
	go a.Server.ListenAndServe()
}
