package usecases

import (
	"context"

	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya-admin/protos"
)

// GetDocumentationForServerType gets documentation for given server type
func GetDocumentationForServerType(serverType string) (map[string]interface{}, error) {
	var ret map[string]interface{}
	responseProto := &protos.Doc{}
	route := serverType + "." + serverType + "remote.docs"

	err := pitaya.RPC(context.Background(), route, responseProto, &protos.Doc{})

	if err != nil {
		return nil, err
	}

	err = pitaya.GetSerializer().Unmarshal([]byte(responseProto.GetDoc()), &ret)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

// GetRemotesDoc get servers remotes docs only
func getRemotesDoc(serverType string) (map[string]interface{}, error) {
	ret, err := GetDocumentationForServerType(serverType)

	if err != nil {
		return nil, err
	}

	rM := ret["remotes"]
	retMap, _ := rM.(map[string]interface{})

	return retMap, nil
}

// GetHandlersDoc get servers handlers docs only
func getHandlersDoc(serverType string) (map[string]interface{}, error) {
	ret, err := GetDocumentationForServerType(serverType)

	if err != nil {
		return nil, err
	}

	rM := ret["handlers"]
	retMap, _ := rM.(map[string]interface{})

	return retMap, nil
}

// GetDocForHandlersOrRemotes gets only remote or handler doc
func GetDocForHandlersOrRemotes(servertype, handlerOrRemote string) (map[string]interface{}, error) {
	if handlerOrRemote == "remote" {
		return getRemotesDoc(servertype)
	}
	return getHandlersDoc(servertype)

}

// GetMethodDoc gets documentation for a method given its route and wether its a remote or a handler
func GetMethodDoc(servertype, methodType, route string) (map[string]interface{}, error) {
	ret, err := GetDocForHandlersOrRemotes(servertype, methodType)

	if err != nil {
		return nil, err
	}

	rM := ret[route]
	retMap, _ := rM.(map[string]interface{})

	return retMap, nil
}
