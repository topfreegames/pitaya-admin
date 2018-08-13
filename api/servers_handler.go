package api

import (
	"encoding/json"
	"net/http"
	"strings"

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
		w.Write([]byte("error: failed to retrieve server details"))
	}

}
