package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/pitaya-admin/api"
	"github.com/topfreegames/pitaya-admin/rpc"
	pt "github.com/topfreegames/pitaya-admin/testing"
	"github.com/topfreegames/pitaya/examples/demo/protos"
)

func TestRPC(t *testing.T) {
	t.Parallel()
	if !pt.IsConf {
		pt.ConfApp()
	}
	a := pt.TestApp
	var (
		response   *httptest.ResponseRecorder
		rpcHandler = api.NewRPCHandler(a)
	)

	tables := map[string]struct {
		request *http.Request
		asserts func(response *httptest.ResponseRecorder)
	}{
		"test_empty_body": {
			request: func() *http.Request {
				request, err := http.NewRequest("POST", "/rpc", nil)
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, response.Code)
				b := response.Body.String()
				assert.Equal(t, b, `{"success":false, "message":"request body shouldnt be empty", "reason": "empty request body"}`)

			},
		},
		"test_bad_post": {
			request: func() *http.Request {
				body := "not a rpc request struct"
				bts, _ := json.Marshal(body)

				request, err := http.NewRequest("POST", "/rpc", bytes.NewReader(bts))
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, response.Code)
				b := response.Body.String()
				assert.Equal(t, b, `{"success":false, "message":"failed to decode request", "reason": "json: cannot unmarshal string into Go value of type rpc.Request"}`)

			},
		},
		"test_bad_route": {
			request: func() *http.Request {
				body := &rpc.Request{
					Route:        "ro.u.te",
					FrontendType: "",
					ServerID:     "123",
					Meta:         "test",
				}
				bts, _ := json.Marshal(body)

				request, err := http.NewRequest("POST", "/rpc", bytes.NewReader(bts))
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, response.Code)
			},
		},
		"test_success_rpc": {
			request: func() *http.Request {
				rpcMessage := &protos.RPCMsg{Msg: "hi im a rpc testing msg"}
				jsonProtobuffMarshaler := jsonpb.Marshaler{}
				rpcMessageSerialized, err := jsonProtobuffMarshaler.MarshalToString(rpcMessage)
				assert.NoError(t, err)
				body := &rpc.Request{
					Route:        "connector.connectorremote.remotefunc",
					FrontendType: "connector",
					ServerID:     "",
					Meta:         rpcMessageSerialized,
				}
				bts, _ := json.Marshal(body)

				request, err := http.NewRequest("POST", "/rpc", bytes.NewReader(bts))
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				b := response.Body.String()
				assert.Equal(t, http.StatusOK, response.Code)
				assert.Equal(t, `{"success" : true, "response":{"Msg":"hi im a rpc testing msg"}}`, b)
			},
		},
		"id_not_correct": {
			request: func() *http.Request {
				rpcMessage := &protos.RPCMsg{Msg: "hi im a rpc msg"}
				jsonProtobuffMarshaler := jsonpb.Marshaler{}
				rpcMessageSerialized, err := jsonProtobuffMarshaler.MarshalToString(rpcMessage)
				assert.NoError(t, err)
				body := &rpc.Request{
					Route:        "connector.connectorremote.remotefunc",
					FrontendType: "connector",
					ServerID:     "wops",
					Meta:         rpcMessageSerialized,
				}
				bts, _ := json.Marshal(body)

				request, err := http.NewRequest("POST", "/rpc", bytes.NewReader(bts))
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, response.Code)
			},
		},
	}

	for name, table := range tables {
		t.Run(name, func(t *testing.T) {
			response = httptest.NewRecorder()
			rpcHandler.ServeHTTP(response, table.request)
			table.asserts(response)
		})
	}
}
