package api

import (
	"encoding/json"
	"net/http"

	"github.com/topfreegames/pitaya/cluster"

	"github.com/topfreegames/pitaya"
)

// ServersHandler handler
type ServersHandler struct {
	App *App
}

// NewServersHandler creates a new servers handler
func NewServersHandler(a *App) *ServersHandler {
	m := &ServersHandler{App: a}
	return m
}

// ServeHTTP method
func (s *ServersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	param := v.Get("type")

	if param == "" {
		serversMap := pitaya.GetServers()

		servers, err := json.Marshal(serversMap)

		if err != nil {
			Write(w, 500, `{"success":false, "message":"failed to marshal response"}`)
			return
		}

		w.Write(servers)

		return
	}

	serversMap, err := pitaya.GetServersByType(param)
	if err != nil {
		Write(w, 500, `{"success":false, "message":"failed to retrieve servers"}`)
		return
	}

	arr := make([]*cluster.Server, 0)
	for _, server := range serversMap {
		arr = append(arr, server)
	}

	servers, err := json.Marshal(arr)

	if err != nil {
		Write(w, 500, `{"success":false, "message":"failed to marshal response"}`)
		return
	}

	w.Write(servers)
}
