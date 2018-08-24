package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya-admin/protos"
)

// DocsHandler handle docs
type DocsHandler struct {
	App *App
}

// NewDocsHandler create new docs handler
func NewDocsHandler(a *App) *DocsHandler {
	m := &DocsHandler{App: a}
	return m
}

func (s *DocsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	param := v.Get("type")
	roh := v.Get("t")
	routeParam := v.Get("route")

	if param != "" {
		rs := &protos.Doc{}
		route := param + "." + param + "remote.docs"
		err := pitaya.RPC(context.Background(), route, rs, &protos.Doc{})

		if err != nil {
			WriteError(w, 500, "failed to send RPC, server must have remote service docs", err)
			return
		}

		if roh != "" {
			var ret map[string]interface{}
			json.Unmarshal([]byte(rs.GetDoc()), &ret)
			rM := ret[roh]
			retMap, _ := rM.(map[string]interface{})
			if routeParam != "" {
				retMapBytes, err := json.Marshal(retMap[routeParam])
				if err != nil {
					WriteError(w, 500, "failed to marshal documentation", err)
					return
				}
				WriteSuccessWithJSON(w, http.StatusOK, retMapBytes, "Successfully got docs")
				return
			}
			retMapBytes, err := json.Marshal(rM)
			if err != nil {
				WriteError(w, 500, "failed to marshal documentation", err)
				return
			}
			WriteSuccessWithJSON(w, http.StatusOK, retMapBytes, "Successfully got docs")
			return

		}
		WriteSuccessWithJSON(w, http.StatusOK, []byte(rs.GetDoc()), "Successfully got docs")
		return
	}

	WriteError(w, http.StatusBadRequest, "must specify at least server type", errors.New("server type not specified"))
}
