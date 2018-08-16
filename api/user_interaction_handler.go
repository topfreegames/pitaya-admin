package api

import (
	"encoding/json"
	"net/http"

	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/modules"
	"github.com/topfreegames/pitaya/protos"
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

type kickMsg struct {
	Uid          string
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
		Write(w, http.StatusBadRequest, `{"success":false, "message":"failed to decode request body into kick struct"}`)
		return
	}

	if s.App.usesGrpc {
		bs, err := pitaya.GetModule("bindingsStorage")
		if err != nil {
			Write(w, 500, `{"success":false, "message":"binding storage must be set when using grpc"}`)
			return
		}

		bindingStorage, ok := bs.(*modules.ETCDBindingStorage)
		if !ok {
			Write(w, 500, `{"success":false, "message":"binding storage must be set when using grpc"}`)
			return
		}

		svID, err := bindingStorage.GetUserFrontendID(kick.Uid, kick.FrontendType)
		if err != nil {
			Write(w, 500, `{"success":false, "message":"error trying to get user frontend server id"}`)
			return
		}

		rpcRoute := kick.FrontendType + ".sys.kick"
		err = pitaya.RPCTo(r.Context(), svID, rpcRoute, &protos.KickAnswer{}, &protos.KickMsg{UserId: kick.Uid})
		if err != nil {
			Write(w, 500, `{"success":false, "message":"failed to send kick rpc"}`)
		}

		Write(w, http.StatusOK, `{"success":true, "message":"user kicked"}`)
		return
	}

	Write(w, 500, `{"success":false, "message":"non-grpc kick not implemented yet"}`)
}
