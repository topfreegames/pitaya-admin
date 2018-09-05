package api

import (
	"encoding/json"
	"net/http"

	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya-admin/constants"
)

// PushMsg post
type PushMsg struct {
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
	var push PushMsg

	if r.Body == nil {
		WriteError(w, http.StatusBadRequest, "request body shouldnt be empty", constants.ErrEmptyRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&push)

	if err != nil {
		WriteError(w, http.StatusBadRequest, "failed to decode request body into push struct", err)
		return
	}

	//calling SendPushToUsers from a backend server, must have frontend type specified
	if push.FrontendType == "" {
		WriteError(w, http.StatusBadRequest, "server type needs to be specified", constants.ErrNoServerType)
		return
	}

	err = pitaya.SendPushToUsers(push.Route, push.Message, push.Uids, push.FrontendType)

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to send push to user", err)
		return
	}

	Write(w, http.StatusOK, `{"success": true}`)

}

// KickMsg post
type KickMsg struct {
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
	var kick KickMsg

	if r.Body == nil {
		WriteError(w, http.StatusBadRequest, "request body shouldnt be empty", constants.ErrEmptyRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&kick)

	if err != nil {
		WriteError(w, http.StatusBadRequest, "failed to decode request body into kick struct", err)
		return
	}

	if kick.FrontendType == "" {
		WriteError(w, http.StatusBadRequest, "server type needs to be specified", constants.ErrNoServerType)
		return
	}

	err = pitaya.SendKickToUsers(kick.Uids, kick.FrontendType)

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to kick user", err)
		return
	}

	Write(w, http.StatusOK, `{"success":true}`)

}
