package api

import (
	"encoding/json"
	"net/http"

	"github.com/topfreegames/pitaya"
)

type pushMsg struct {
	Uids         []string
	Route        string
	Message      interface{}
	FrontendType string
}

// PushToUsersHandler handler
type PushToUsersHandler struct {
	App *App
}

// NewPushToUsersHandler creates a new push to users handler
func NewPushToUsersHandler(a *App) *PushToUsersHandler {
	m := &PushToUsersHandler{App: a}
	return m
}

// ServeHTTP method
func (s *PushToUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var push pushMsg
	err := json.NewDecoder(r.Body).Decode(&push)

	if err != nil {
		Write(w, http.StatusBadRequest, `{"success":false, "message":"failed to decode request body into push struct"}`)
		return
	}

	err = pitaya.SendPushToUsers(push.Route, push.Message, push.Uids, push.FrontendType)

	if err != nil {
		Write(w, 500, `{"success":false, "message":"error trying to push message to user"}`)
		return
	}

	Write(w, http.StatusOK, `{"success": true}`)

}
