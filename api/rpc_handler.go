package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/topfreegames/pitaya-admin/usecases"

	"github.com/topfreegames/pitaya"
)

// RPCHandler handler
type RPCHandler struct {
	App *App
}

// NewRPCHandler creates a new handler
func NewRPCHandler(a *App) *RPCHandler {
	m := &RPCHandler{App: a}
	return m
}

// RPCRequest struct
type RPCRequest struct {
	Route        string
	FrontendType string
	ServerID     string
	Meta         string // Data itself
	ProtoSource  []byte // This is the proto source file for protos that the method in route will use
}

func (s *RPCHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var remote RPCRequest
	err := json.NewDecoder(r.Body).Decode(&remote)

	requestMessage, responseMessage, err := usecases.CreateRPCMessagesFromProto(remote.FrontendType, remote.Route, remote.ProtoSource, remote.Meta)

	if err != nil {
		WriteError(w, http.StatusBadRequest, "failed to create RPC", err)
		return
	}

	if remote.ServerID != "" {
		err = pitaya.RPCTo(context.Background(), remote.ServerID, remote.Route, responseMessage, requestMessage)
	} else {
		err = pitaya.RPC(context.Background(), remote.Route, responseMessage, requestMessage)
	}

	if err != nil {
		WriteError(w, 500, "failed to send RPC", err)
		return
	}

	jsonResponse, err := responseMessage.MarshalJSON()

	if err != nil {
		WriteError(w, 500, "failed to marshal response into JSON", err)
		return
	}

	WriteSuccessWithJSON(w, http.StatusOK, jsonResponse, "Successfully sent RPC")
}
