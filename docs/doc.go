package docs

import (
	"context"

	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya-admin/protos"
)

// GetDocumentationForServerType gets documentation for given server type.
func GetDocumentationForServerType(serverType, docsRemoteRoute string, flag bool) (map[string]interface{}, error) {
	var ret map[string]interface{}
	responseProto := &protos.Doc{}
	route := serverType + "." + serverType + docsRemoteRoute

	err := pitaya.RPC(context.Background(), route, responseProto, &protos.DocMsg{GetProtos: flag})

	if err != nil {
		return nil, err
	}

	err = pitaya.GetSerializer().Unmarshal([]byte(responseProto.GetDoc()), &ret)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func getRemotesDoc(serverType, docsRemoteRoute string, flag bool) (map[string]interface{}, error) {
	ret, err := GetDocumentationForServerType(serverType, docsRemoteRoute, flag)

	if err != nil {
		return nil, err
	}

	rM := ret["remotes"]
	retMap, _ := rM.(map[string]interface{})

	return retMap, nil
}

func getHandlersDoc(serverType, docsRemoteRoute string, flag bool) (map[string]interface{}, error) {
	ret, err := GetDocumentationForServerType(serverType, docsRemoteRoute, flag)

	if err != nil {
		return nil, err
	}

	rM := ret["handlers"]
	retMap, _ := rM.(map[string]interface{})

	return retMap, nil
}

// GetDocForHandlersOrRemotes gets only remote or handler docs for given servertype
func GetDocForHandlersOrRemotes(servertype, handlerOrRemote, docsRemoteRoute string, flag bool) (map[string]interface{}, error) {
	if handlerOrRemote == "remote" {
		return getRemotesDoc(servertype, docsRemoteRoute, flag)
	}
	return getHandlersDoc(servertype, docsRemoteRoute, flag)

}

// GetMethodDoc gets documentation for a method given its route and whether it's a remote or a handler
func GetMethodDoc(servertype, methodType, route, docsRemoteRoute string, flag bool) (map[string]interface{}, error) {
	ret, err := GetDocForHandlersOrRemotes(servertype, methodType, docsRemoteRoute, flag)

	if err != nil {
		return nil, err
	}

	rM := ret[route]
	retMap, _ := rM.(map[string]interface{})

	return retMap, nil
}
