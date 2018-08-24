package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/topfreegames/pitaya"

	"github.com/jhump/protoreflect/dynamic"

	"github.com/jhump/protoreflect/desc/protoparse"
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
	Meta         string   // Data itself
	ProtoSource  []byte   // This is the proto source file for protos that the method in route will use
	ArgsName     []string // This is input and output proto names, ie, RPCMes, RPCAns, etc
}

func (s *RPCHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var remote RPCRequest
	err := json.NewDecoder(r.Body).Decode(&remote)

	if err != nil {
		WriteError(w, http.StatusBadRequest, "failed to unmarshal documentation", err)
		return
	}

	parserAccessor := func(filename string) (io.ReadCloser, error) {
		return ioutil.NopCloser(bytes.NewReader(remote.ProtoSource)), nil
	}
	parser := protoparse.Parser{Accessor: parserAccessor}
	fileDescriptor, err := parser.ParseFiles("")

	if err != nil {
		WriteError(w, http.StatusBadRequest, "failed to parse proto file", err)
		return
	}

	requestMessage := dynamic.NewMessage(fileDescriptor[0].FindMessage(remote.ArgsName[0]))
	responseMessage := dynamic.NewMessage(fileDescriptor[0].FindMessage(remote.ArgsName[1]))
	err = requestMessage.UnmarshalJSON([]byte(remote.Meta)) //test if this can come as byte alrdy

	if err != nil {
		WriteError(w, http.StatusBadRequest, "failed to unmarshal proto message", err)
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

	WriteSuccessWithJSON(w, http.StatusOK, jsonResponse)
}
