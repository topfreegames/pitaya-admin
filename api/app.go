package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/topfreegames/pitaya"
)

type app struct {
	router *mux.Router
	host   string
	port   int
}

type pushMsg struct {
	Uids         []string
	Route        string
	Message      interface{}
	FrontendType string
}

func getServersByType(w http.ResponseWriter, r *http.Request) {
	param := strings.Replace(r.URL.Path, "/servers/", "", -1)
	serversMap, err := pitaya.GetServersByType(param)

	if err != nil {
		w.Write([]byte("error: failed to retrieve servers by type"))
		return
	}

	for _, server := range serversMap {
		serverJSON, err := json.Marshal(server)
		if err == nil {
			w.Write(serverJSON)
			continue
		}
		w.Write([]byte("error: failed to retrieve server specs"))
	}
}

func sendPushToUser(w http.ResponseWriter, r *http.Request) {
	var push pushMsg
	err := json.NewDecoder(r.Body).Decode(&push)
	if err == nil {
		err = pitaya.SendPushToUsers(push.Route, push.Message, push.Uids, push.FrontendType)
		if err != nil {
			w.Write([]byte("error: failed to push message to users"))
		} else {
			w.Write([]byte("message pushed to users"))
		}
	}
}

func newRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/servers/{type}", getServersByType).Methods("GET")
	router.HandleFunc("/user/push", sendPushToUser).Methods("POST")
	return router

}

// Init starts the api
func Init() {
	api := &app{
		router: newRouter(),
		host:   "localhost",
		port:   8000,
	}
	http.ListenAndServe(fmt.Sprintf("%s:%d", api.host, api.port), api.router)
}
