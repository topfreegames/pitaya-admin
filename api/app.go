package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// App struct
type App struct {
	Router *mux.Router
	host   string
	port   int
}

func (a *App) getRouter() *mux.Router {
	router := mux.NewRouter()
	router.Handle("/servers/{type}", NewServersHandler(a)).Methods("GET")
	router.Handle("/user/push", NewPushToUsersHandler(a)).Methods("POST")
	return router
}

// Init starts the api
func Init() {
	api := &App{
		host: "localhost",
		port: 8000,
	}
	api.Router = api.getRouter()
	http.ListenAndServe(fmt.Sprintf("%s:%d", api.host, api.port), api.Router)
}
