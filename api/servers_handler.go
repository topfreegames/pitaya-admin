package api

import (
	"encoding/json"
	"net/http"

	"github.com/topfreegames/pitaya/cluster"

	"github.com/topfreegames/pitaya"
)

// ServersHandler handle list server routes
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
			WriteError(w, http.StatusInternalServerError, "failed trying to marshal servers map", err)
			return
		}

		WriteSuccessWithJSON(w, http.StatusOK, servers, "Successfully got servers")
		return
	}

	serversMap, err := pitaya.GetServersByType(param)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to retrieve servers", err)
		return
	}

	arr := make([]*cluster.Server, 0)
	for _, server := range serversMap {
		arr = append(arr, server)
	}

	servers, err := json.Marshal(arr)

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to marshal servers array", err)
		return
	}

	WriteSuccessWithJSON(w, http.StatusOK, servers, "Successfully got servers")
}
