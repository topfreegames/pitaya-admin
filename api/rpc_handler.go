package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya-admin/constants"
	"github.com/topfreegames/pitaya-admin/rpc"
	pconstants "github.com/topfreegames/pitaya/constants"
)

// RPCHandler handles RPC routes
type RPCHandler struct {
	App *App
}

// NewRPCHandler creates a new handler
func NewRPCHandler(a *App) *RPCHandler {
	m := &RPCHandler{App: a}
	return m
}

// ServeHTTP method
func (s *RPCHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var remote rpc.Request

	if r.Body == nil {
		WriteError(w, http.StatusBadRequest, "request body shouldnt be empty", constants.ErrEmptyRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&remote)

	if err != nil {
		WriteError(w, http.StatusBadRequest, "failed to decode request", err)
		return
	}

	requestMessage, responseMessage, err := rpc.CreateRPCMessagesFromProto(remote, s.App.GetRemoteDocsRoute(), s.App.GetRemoteProtosRoute())

	if err != nil {
		WriteError(w, http.StatusBadRequest, "failed to create RPC", err)
		return
	}

	if remote.ServerID != "" {
		err = pitaya.RPCTo(context.Background(), remote.ServerID, remote.Route, responseMessage, requestMessage)
		//checks if error was caused by not found ID to return 404, otherwise return 500 below
		if err != nil && err == pconstants.ErrServerNotFound {
			WriteError(w, http.StatusNotFound, "server not found", err)
			return
		}

	} else {
		err = pitaya.RPC(context.Background(), remote.Route, responseMessage, requestMessage)
	}

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to send RPC", err)
		return
	}

	jsonResponse, err := responseMessage.MarshalJSON()

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to marshal response into JSON", err)
		return
	}

	WriteSuccessWithJSON(w, http.StatusOK, jsonResponse, "Successfully sent RPC")
}
