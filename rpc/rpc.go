package rpc

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"

	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/topfreegames/pitaya-admin/constants"
	"github.com/topfreegames/pitaya-admin/docs"
)

// RPCRequest represents a request received by the server to be used in RPC
type RPCRequest struct {
	Route        string
	FrontendType string
	ServerID     string
	Meta         string // Data itself
	ProtoSource  []byte // This is the proto source file for protos that the method in route will use
}

func getInputAndOuputProtosFromAutodoc(route, frontendtype, docsRemoteRoute string) (string, string, error) {

	outputName := ""
	inputName := ""

	remoteMethodDocs, err := docs.GetMethodDoc(frontendtype, "remote", route, docsRemoteRoute)

	if err != nil {
		return "", "", err
	}

	in := remoteMethodDocs["input"]
	inputDocs, ok := in.(map[string]interface{})

	if !ok {
		return "", "", constants.ErrNoInputDocForMethod
	}

	out := remoteMethodDocs["output"]
	outputDocsArr := out.([]interface{})
	outputDocs, ok := outputDocsArr[0].(map[string]interface{})

	if !ok {
		return "", "", constants.ErrNoOutputDocForMethod
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
		return "", "", constants.ErrProtoDoc
	}

	return inputName, outputName, nil

}

// CreateRPCMessagesFromProto outputs the reply and args messages types for a given RPC route using protoparse and dynamic messages
func CreateRPCMessagesFromProto(frontendtype, route string, protosource []byte, meta, docsRemoteRoute string) (*dynamic.Message, *dynamic.Message, error) {

	inputName, outputName, err := getInputAndOuputProtosFromAutodoc(route, frontendtype, docsRemoteRoute)

	if err != nil {
		return nil, nil, err
	}

	parserAccessor := func(filename string) (io.ReadCloser, error) {
		return ioutil.NopCloser(bytes.NewReader(protosource)), nil
	}
	parser := protoparse.Parser{Accessor: parserAccessor}
	fileDescriptor, err := parser.ParseFiles("")

	if err != nil {
		return nil, nil, err
	}

	requestMessage := dynamic.NewMessage(fileDescriptor[0].FindMessage(inputName))
	responseMessage := dynamic.NewMessage(fileDescriptor[0].FindMessage(outputName))
	err = requestMessage.UnmarshalJSON([]byte(meta))

	return requestMessage, responseMessage, nil
}
