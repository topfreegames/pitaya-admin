package rpc

import (
	"bytes"
	"compress/gzip"
	"context"
	"io/ioutil"
	"strings"

	"github.com/gogo/protobuf/proto"
	protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya-admin/constants"
	"github.com/topfreegames/pitaya-admin/docs"
	"github.com/topfreegames/pitaya-admin/protos"
)

// Request represents a request received by the server to be used in RPC
type Request struct {
	Route        string
	FrontendType string
	ServerID     string
	Meta         string // Data itself
}

func getInputAndOuputProtosFromAutodoc(route, frontendtype, docsRemoteRoute string) (string, string, error) {

	outputName := ""
	inputName := ""

	remoteMethodDocs, err := docs.GetMethodDoc(frontendtype, "remote", route, docsRemoteRoute, true)

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

func unpackDescriptor(compressedDescriptor []byte) (*protobuf.FileDescriptorProto, error) {
	r, err := gzip.NewReader(bytes.NewReader(compressedDescriptor))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, constants.ErrGzip
	}

	fileDescriptorProto := new(protobuf.FileDescriptorProto)

	err = proto.Unmarshal(b, fileDescriptorProto)
	if err != nil {
		return nil, err
	}

	return fileDescriptorProto, nil
}

func retrieveDescriptor(protoName, protosRemoteRoute string) (*desc.FileDescriptor, error) {
	replyMsg := &protos.ProtoDescriptor{}
	protoRequest := &protos.ProtoName{Name: protoName}

	err := pitaya.RPC(context.Background(), protosRemoteRoute, replyMsg, protoRequest)

	if err != nil {
		return nil, err
	}

	fileDescriptorProto, err := unpackDescriptor(replyMsg.Desc)

	if err != nil {
		return nil, err
	}

	depsArr := []*desc.FileDescriptor{}

	for _, dep := range fileDescriptorProto.GetDependency() {
		replyMsgForDeps := &protos.ProtoDescriptor{}
		protoRequestForDeps := &protos.ProtoName{Name: dep}
		err = pitaya.RPC(context.Background(), protosRemoteRoute, replyMsgForDeps, protoRequestForDeps)

		if err != nil {
			return nil, err
		}

		depDescriptorProto, err := unpackDescriptor(replyMsgForDeps.Desc)

		if err != nil {
			return nil, err
		}

		transitiveDepArr := []*desc.FileDescriptor{}

		for _, transitiveDep := range depDescriptorProto.GetDependency() {
			transitiveDepDescriptor, err := retrieveDescriptor(transitiveDep, protosRemoteRoute)
			if err != nil {
				return nil, err
			}
			transitiveDepArr = append(transitiveDepArr, transitiveDepDescriptor)
		}

		depDescriptor, err := desc.CreateFileDescriptor(depDescriptorProto, transitiveDepArr...)

		if err != nil {
			return nil, err
		}

		depsArr = append(depsArr, depDescriptor)
	}

	fileDescriptor, err := desc.CreateFileDescriptor(fileDescriptorProto, depsArr...)
	return fileDescriptor, err
}

func buildMessage(serverType, protoName, protosRemoteRoute string) (*dynamic.Message, error) {
	rpcRoute := serverType + "." + protosRemoteRoute
	fileDescriptor, err := retrieveDescriptor(protoName, rpcRoute)

	if err != nil {
		return nil, err
	}

	return dynamic.NewMessage(fileDescriptor.FindMessage(protoName)), nil
}

// CreateRPCMessagesFromProto outputs the reply and args messages types for a given RPC route using protoparse and dynamic messages
func CreateRPCMessagesFromProto(request Request, docsRemoteRoute, protosRemoteRoute string) (*dynamic.Message, *dynamic.Message, error) {

	inputName, outputName, err := getInputAndOuputProtosFromAutodoc(request.Route, request.FrontendType, docsRemoteRoute)

	if err != nil {
		return nil, nil, err
	}

	requestMessage, err := buildMessage(request.FrontendType, inputName, protosRemoteRoute)

	if err != nil {
		return nil, nil, err
	}

	responseMessage, err := buildMessage(request.FrontendType, outputName, protosRemoteRoute)

	if err != nil {
		return nil, nil, err
	}

	err = requestMessage.UnmarshalJSON([]byte(request.Meta))

	if err != nil {
		return nil, nil, err
	}

	return requestMessage, responseMessage, nil
}
