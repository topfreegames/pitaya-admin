package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/topfreegames/pitaya/cluster"
	"github.com/topfreegames/pitaya/logger"

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
			logger.Log.Errorf("error trying to marshal servers map: %s", err.Error())
			errMsg := fmt.Sprintf(`{"success":false, "message":"failed to marshal response", "reason": "%s"}`, err.Error())
			Write(w, 500, errMsg)
			return
		}

		w.Write(servers)

		return
	}

	serversMap, err := pitaya.GetServersByType(param)
	if err != nil {
		logger.Log.Errorf("error trying to retrieve servers: %s", err.Error())
		errMsg := fmt.Sprintf(`{"success":false, "message":"failed to retrieve servers", "reason": "%s"}`, err.Error())
		Write(w, 500, errMsg)
		return
	}

	arr := make([]*cluster.Server, 0)
	for _, server := range serversMap {
		arr = append(arr, server)
	}

	servers, err := json.Marshal(arr)

	if err != nil {
		logger.Log.Errorf("error trying to marshal servers array: %s", err.Error())
		errMsg := fmt.Sprintf(`{"success":false, "message":"failed to marshal response", "reason": "%s"}`, err.Error())
		Write(w, 500, errMsg)
		return
	}

	w.Write(servers)
}
