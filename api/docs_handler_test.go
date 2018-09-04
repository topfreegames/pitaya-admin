package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/pitaya-admin/api"
	pt "github.com/topfreegames/pitaya-admin/testing"
)

func TestGetDocs(t *testing.T) {
	t.Parallel()
	if pt.IsConf == false {
		pt.ConfApp()
	}
	a := pt.TestApp
	var (
		response    *httptest.ResponseRecorder
		docsHandler *api.DocsHandler
	)
	docsHandler = api.NewDocsHandler(a)
	tables := map[string]struct {
		request *http.Request
		asserts func(response *httptest.ResponseRecorder)
	}{
		"test_get_docs_no_param": {
			request: func() *http.Request {
				request, err := http.NewRequest("GET", "/docs", nil)
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				b := response.Body.String()
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Equal(t, `{"success":false, "message":"at least server type must be specified", "reason": "no server type chosen for getting autodoc"}`, b)
			},
		},
		"test_get_docs_by_type_but_sv_does_not_exists": {
			request: func() *http.Request {
				request, err := http.NewRequest("GET", "/docs?type=bla", nil)
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				b := response.Body.String()
				assert.Equal(t, http.StatusInternalServerError, response.Code)
				assert.Equal(t, `{"success":false, "message":"failed to get docs", "reason": "no servers available of this type"}`, b)
			},
		},
		"test_get_docs_success": {
			request: func() *http.Request {
				request, err := http.NewRequest("GET", "/docs?type=connector", nil)
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				b := response.Body.String()
				assert.Equal(t, `{"handlers":{"connector.connector.getsessiondata":{"input":null,"output":[{"Data":"map[string]interface {}"},"error"]},"connector.connector.notifysessiondata":{"input":{"data":"map[string]interface {}"},"output":[]},"connector.connector.setsessiondata":{"input":{"data":"map[string]interface {}"},"output":[{"Code":"int32","Msg":"string"},"error"]}},"remotes":{"connector.connectorremote.docs":{"input":{"getProtos":"bool"},"output":[{"doc":"string"},"error"]},"connector.connectorremote.proto":{"input":{"name":"string"},"output":[{"desc":"[]byte"},"error"]},"connector.connectorremote.remotefunc":{"input":{"Msg":"string"},"output":[{"Msg":"string"},"error"]},"connector.sys.bindsession":{"input":{"data":"[]byte","id":"int64","uid":"string"},"output":[{"data":"[]byte","error":{"code":"string","metadata":"map[string]string","msg":"string"}},"error"]},"connector.sys.kick":{"input":{"userId":"string"},"output":[{"kicked":"bool"},"error"]},"connector.sys.pushsession":{"input":{"data":"[]byte","id":"int64","uid":"string"},"output":[{"data":"[]byte","error":{"code":"string","metadata":"map[string]string","msg":"string"}},"error"]}}}`, b)
				assert.Equal(t, http.StatusOK, response.Code)
			},
		},
		"test_get_docs_success_with_protos": {
			request: func() *http.Request {
				request, err := http.NewRequest("GET", "/docs?type=connector&getProtos=1", nil)
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				b := response.Body.String()
				assert.Equal(t, `{"handlers":{"connector.connector.getsessiondata":{"input":null,"output":[{"*services.SessionData":{"Data":"map[string]interface {}"}},"error"]},"connector.connector.notifysessiondata":{"input":{"*services.SessionData":{"data":"map[string]interface {}"}},"output":[]},"connector.connector.setsessiondata":{"input":{"*services.SessionData":{"data":"map[string]interface {}"}},"output":[{"*services.Response":{"Code":"int32","Msg":"string"}},"error"]}},"remotes":{"connector.connectorremote.docs":{"input":{"*protos.DocMsg":{"getProtos":"bool"}},"output":[{"*protos.Doc":{"doc":"string"}},"error"]},"connector.connectorremote.proto":{"input":{"*protos.ProtoName":{"name":"string"}},"output":[{"*protos.ProtoDescriptor":{"desc":"[]byte"}},"error"]},"connector.connectorremote.remotefunc":{"input":{"*protos.RPCMsg":{"Msg":"string"}},"output":[{"*protos.RPCRes":{"Msg":"string"}},"error"]},"connector.sys.bindsession":{"input":{"*protos.Session":{"data":"[]byte","id":"int64","uid":"string"}},"output":[{"*protos.Response":{"data":"[]byte","error":{"*protos.Error":{"code":"string","metadata":"map[string]string","msg":"string"}}}},"error"]},"connector.sys.kick":{"input":{"*protos.KickMsg":{"userId":"string"}},"output":[{"*protos.KickAnswer":{"kicked":"bool"}},"error"]},"connector.sys.pushsession":{"input":{"*protos.Session":{"data":"[]byte","id":"int64","uid":"string"}},"output":[{"*protos.Response":{"data":"[]byte","error":{"*protos.Error":{"code":"string","metadata":"map[string]string","msg":"string"}}}},"error"]}}}`, b)
				assert.Equal(t, http.StatusOK, response.Code)
			},
		},
		"test_success_with_params": {
			request: func() *http.Request {
				request, err := http.NewRequest("GET", "/docs?type=connector&methodtype=remote", nil)
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				b := response.Body.String()
				assert.Equal(t, `{"connector.connectorremote.docs":{"input":{"getProtos":"bool"},"output":[{"doc":"string"},"error"]},"connector.connectorremote.proto":{"input":{"name":"string"},"output":[{"desc":"[]byte"},"error"]},"connector.connectorremote.remotefunc":{"input":{"Msg":"string"},"output":[{"Msg":"string"},"error"]},"connector.sys.bindsession":{"input":{"data":"[]byte","id":"int64","uid":"string"},"output":[{"data":"[]byte","error":{"code":"string","metadata":"map[string]string","msg":"string"}},"error"]},"connector.sys.kick":{"input":{"userId":"string"},"output":[{"kicked":"bool"},"error"]},"connector.sys.pushsession":{"input":{"data":"[]byte","id":"int64","uid":"string"},"output":[{"data":"[]byte","error":{"code":"string","metadata":"map[string]string","msg":"string"}},"error"]}}`, b)
				assert.Equal(t, http.StatusOK, response.Code)
			},
		},
		"test_fail_with_params": {
			request: func() *http.Request {
				request, err := http.NewRequest("GET", "/docs?type=fail&methodtype=remote", nil)
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, response.Code)
			},
		},
		"test_success_with_full_params": {
			request: func() *http.Request {
				request, err := http.NewRequest("GET", "/docs?type=connector&methodtype=remote&route=connector.connectorremote.docs", nil)
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				b := response.Body.String()
				assert.Equal(t, `{"input":{"getProtos":"bool"},"output":[{"doc":"string"},"error"]}`, b)
				assert.Equal(t, http.StatusOK, response.Code)
			},
		},
		"test_fail_with_full_params": {
			request: func() *http.Request {
				request, err := http.NewRequest("GET", "/docs?type=fail&methodtype=remote&route=connector.connectorremote.docs", nil)
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, response.Code)
			},
		},
	}

	for name, table := range tables {
		t.Run(name, func(t *testing.T) {
			response = httptest.NewRecorder()
			docsHandler.ServeHTTP(response, table.request)
			table.asserts(response)
		})
	}
}
