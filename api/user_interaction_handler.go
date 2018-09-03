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

// PushToUsersHandler handle push route
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
		WriteError(w, http.StatusInternalServerError, "failed to decode request body into push struct", err)
		return
	}

	err = pitaya.SendPushToUsers(push.Route, push.Message, push.Uids, push.FrontendType)

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to send push to user", err)
		return
	}

	Write(w, http.StatusOK, `{"success": true}`)

}

type kickMsg struct {
	Uids         []string
	FrontendType string
}

// KickUserHandler handle kick route
type KickUserHandler struct {
	App *App
}

// NewKickUserHandler creates a new kick user handler
func NewKickUserHandler(a *App) *KickUserHandler {
	m := &KickUserHandler{App: a}
	return m
}

func (s *KickUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var kick kickMsg
	err := json.NewDecoder(r.Body).Decode(&kick)

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to decode request body into kick struct", err)
		return
	}

	err = pitaya.SendKickToUsers(kick.Uids, kick.FrontendType)

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to kick user", err)
		return
	}

	Write(w, http.StatusOK, `{"success":true}`)

}
