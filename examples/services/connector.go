package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya-admin/examples/protos"
	pitayaprotos "github.com/topfreegames/pitaya-admin/protos"
	"github.com/topfreegames/pitaya/component"
)

// ConnectorRemote is a remote that will receive rpc's
type ConnectorRemote struct {
	component.Base
}

// Connector struct
type Connector struct {
	component.Base
}

// SessionData struct
type SessionData struct {
	Data map[string]interface{}
}

// Response struct
type Response struct {
	Code int32
	Msg  string
}

func reply(code int32, msg string) (*Response, error) {
	res := &Response{
		Code: code,
		Msg:  msg,
	}
	return res, nil
}

// GetSessionData gets the session data
func (c *Connector) GetSessionData(ctx context.Context) (*SessionData, error) {
	s := pitaya.GetSessionFromCtx(ctx)
	res := &SessionData{
		Data: s.GetData(),
	}
	return res, nil
}

// SetSessionData sets the session data
func (c *Connector) SetSessionData(ctx context.Context, data *SessionData) (*Response, error) {
	s := pitaya.GetSessionFromCtx(ctx)
	err := s.SetData(data.Data)
	if err != nil {
		return nil, pitaya.Error(err, "CN-000", map[string]string{"failed": "set data"})
	}
	return reply(200, "success")
}

// NotifySessionData sets the session data
func (c *Connector) NotifySessionData(ctx context.Context, data *SessionData) {
	s := pitaya.GetSessionFromCtx(ctx)
	err := s.SetData(data.Data)
	if err != nil {
		fmt.Println("got error on notify", err)
	}
}

// RemoteFunc is a function that will be called remotely
func (c *ConnectorRemote) RemoteFunc(ctx context.Context, msg *protos.RPCMsg) (*protos.RPCRes, error) {
	return &protos.RPCRes{
		Msg: msg.GetMsg(),
	}, nil
}

// Proto is a function that will be called remotely to get the game proto specified by name.doc
func (c *ConnectorRemote) Proto(ctx context.Context, name *pitayaprotos.ProtoName) (*pitayaprotos.ProtoDescriptor, error) {
	protoName := name.Name
	protoReflectTypePointer := proto.MessageType(protoName)
	protoReflectType := protoReflectTypePointer.Elem()
	protoValue := reflect.New(protoReflectType)
	descriptorMethod, ok := protoReflectTypePointer.MethodByName("Descriptor")

	if !ok {
		return nil, errors.New("failed to get proto descriptor")
	}

	descriptorValue := descriptorMethod.Func.Call([]reflect.Value{protoValue})
	protoDescriptor := descriptorValue[0].Bytes()
	return &pitayaprotos.ProtoDescriptor{
		Desc: protoDescriptor,
	}, nil
}

// Docs returns documentation
func (c *ConnectorRemote) Docs(ctx context.Context, flag *pitayaprotos.DocMsg) (*pitayaprotos.Doc, error) {
	d, err := pitaya.Documentation(flag.GetGetProtos())

	if err != nil {
		return nil, err
	}
	doc, err := json.Marshal(d)

	if err != nil {
		return nil, err
	}

	return &pitayaprotos.Doc{Doc: string(doc)}, nil
}
