package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/logger"
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
		logger.Log.Errorf("error trying to decode request body into push struct: %s", err.Error())
		errMsg := fmt.Sprintf(`{"success":false, "message":"failed to decode request body into push struct", "reason": "%s"}`, err.Error())
		Write(w, http.StatusBadRequest, errMsg)
		return
	}

	err = pitaya.SendPushToUsers(push.Route, push.Message, push.Uids, push.FrontendType)

	if err != nil {
		logger.Log.Errorf("error trying to send push to user: %s", err.Error())
		errMsg := fmt.Sprintf(`{"success":false, "message":"failed to send push to user", "reason": "%s"}`, err.Error())
		Write(w, http.StatusBadRequest, errMsg)
		return
	}

	Write(w, http.StatusOK, `{"success": true}`)

}

type kickMsg struct {
	Uids         []string
	FrontendType string
}

// KickUserHandler handler
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
		logger.Log.Errorf("error trying to decode request body into kick struct: %s", err.Error())
		errMsg := fmt.Sprintf(`{"success":false, "message":"failed to decode request body into kick struct", "reason": "%s"}`, err.Error())
		Write(w, http.StatusBadRequest, errMsg)
		return
	}

	err = pitaya.SendKickToUsers(kick.Uids, kick.FrontendType)

	if err != nil {
		logger.Log.Errorf("error trying to kick user: %s", err.Error())
		errMsg := fmt.Sprintf(`{"success":false, "message":"failed to kick user", "reason": "%s"}`, err.Error())
		Write(w, 500, errMsg)
		return
	}

	Write(w, http.StatusOK, `{"success":true}`)

}
