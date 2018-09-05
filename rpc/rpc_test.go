package rpc_test

import (
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/pitaya-admin/constants"
	"github.com/topfreegames/pitaya-admin/rpc"
	pt "github.com/topfreegames/pitaya-admin/testing"
	"github.com/topfreegames/pitaya/examples/demo/protos"
)

func TestCreateRPCMessagesFromProto(t *testing.T) {
	t.Parallel()
	if !pt.IsConf {
		pt.ConfApp()
	}

	rpcMessage := &protos.RPCMsg{Msg: "hi im a rpc msg"}
	jsonProtobuffMarshaler := jsonpb.Marshaler{}
	rpcMessageSerialized, err := jsonProtobuffMarshaler.MarshalToString(rpcMessage)

	assert.NoError(t, err)

	r := rpc.Request{
		Route:        "connector.connectorremote.remotefunc",
		FrontendType: "connector",
		ServerID:     "",
		Meta:         rpcMessageSerialized,
	}

	r2 := rpc.Request{
		Route:        "connector.connectorremote.notAFunc",
		FrontendType: "connector",
		ServerID:     "",
		Meta:         rpcMessageSerialized,
	}

	tables := []struct {
		name              string
		request           rpc.Request
		docsRemoteRoute   string
		protosRemoteRoute string
		err               error
	}{
		{"success", r, "remote.docs", "remote.proto", nil},
		{"fail no remote", r2, "remote.docs", "remote.proto", constants.ErrNoInputDocForMethod},
	}

	for _, table := range tables {
		m1, m2, err := rpc.CreateRPCMessagesFromProto(table.request, table.docsRemoteRoute, table.protosRemoteRoute)
		if table.err == nil {
			assert.NoError(t, err)
			assert.NotNil(t, m1)
			assert.NotNil(t, m2)
		} else {
			assert.Equal(t, table.err, err)
		}
	}
}