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
	if !pt.IsConf {
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
				assert.Equal(t, http.StatusNotFound, response.Code)
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
				assert.Contains(t, b, `"handlers"`)
				assert.Contains(t, b, `"remotes"`)
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
				assert.Contains(t, b, `"handlers"`)
				assert.Contains(t, b, `"remotes"`)
				assert.Contains(t, b, "protos")
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
				assert.Contains(t, b, "connector.sys.bindsession")
				assert.Contains(t, b, "input")
				assert.Contains(t, b, "output")
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
				assert.Contains(t, b, `"input"`)
				assert.Contains(t, b, `"output"`)
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
