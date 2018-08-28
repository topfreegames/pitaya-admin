package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya-admin/protos"

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
	Meta         string // Data itself
	ProtoSource  []byte // This is the proto source file for protos that the method in route will use
}

func (s *RPCHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var remote RPCRequest
	var inputName, outputName string
	inputName = ""
	outputName = ""
	err := json.NewDecoder(r.Body).Decode(&remote)

	if err != nil {
		WriteError(w, http.StatusBadRequest, "failed to unmarshal RPC request", err)
		return
	}

	documentationProto := &protos.Doc{}
	documentationRoute := remote.FrontendType + "." + remote.FrontendType + "remote.docs" // todo- set a default for this sufix
	err = pitaya.RPC(context.Background(), documentationRoute, documentationProto, &protos.Doc{})

	if err != nil {
		WriteError(w, http.StatusBadRequest, "failed to get documentation via RPC", err)
		return
	}

	// Documentation is returned as a map[string]interface, where interfaces are mostly others map[string]interface, making some type assertion`s needed
	var ret map[string]interface{}
	json.Unmarshal([]byte(documentationProto.GetDoc()), &ret)
	rM := ret["remotes"]
	remotesMap := rM.(map[string]interface{})
	rMD := remotesMap[remote.Route]
	remoteMethodDocs := rMD.(map[string]interface{})
	in := remoteMethodDocs["input"]
	inputDocs, ok := in.(map[string]interface{})

	if !ok {
		WriteError(w, 500, "failed to get input documentation for given method", errors.New("failed to retrieve input"))
		return
	}

	out := remoteMethodDocs["output"]
	outputDocsArr := out.([]interface{})
	outputDocs, ok := outputDocsArr[0].(map[string]interface{})

	if !ok {
		WriteError(w, 500, "failed to get output documentation for given method", errors.New("failed to retrieve output"))
		return
	}

	for k := range outputDocs {
		if strings.Contains(k, "proto") {
			outputName = strings.Replace(k, "*", "", 1)
		}
	}

	for k := range inputDocs {
		if strings.Contains(k, "proto") {
			inputName = strings.Replace(k, "*", "", 1)
		}
	}

	if outputName == "" || inputName == "" {
		WriteError(w, 500, "failed to get protos type from docs", errors.New("failed to retrieve protos names"))
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

	requestMessage := dynamic.NewMessage(fileDescriptor[0].FindMessage(inputName))
	responseMessage := dynamic.NewMessage(fileDescriptor[0].FindMessage(outputName))
	err = requestMessage.UnmarshalJSON([]byte(remote.Meta))

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

	WriteSuccessWithJSON(w, http.StatusOK, jsonResponse, "Successfully sent RPC")
}
