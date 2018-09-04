package api_test

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/pitaya-admin/api"
	pt "github.com/topfreegames/pitaya-admin/testing"
)

func TestGetServerList(t *testing.T) {
	t.Parallel()
	if !pt.IsConf {
		pt.ConfApp()
	}
	a := pt.TestApp
	var (
		response      *httptest.ResponseRecorder
		serverHandler = api.NewServersHandler(a)
	)

	tables := map[string]struct {
		request *http.Request
		asserts func(response *httptest.ResponseRecorder)
	}{
		"test_get_servers": {
			request: func() *http.Request {
				request, err := http.NewRequest("GET", "", nil)
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				b := response.Body.String()
				assert.Equal(t, http.StatusOK, response.Code)
				assert.Regexp(t, regexp.MustCompile(`\A\{"success" : true, "response":(\[\{"id":"([a-z0-9-]*)*","type":"([a-z0-9-._]*)","metadata":\{(.)*\},"frontend":(true|false)\},*)+\]\}`), b)
			},
		},
		"test_get_servers_type": {
			request: func() *http.Request {
				request, err := http.NewRequest("GET", "/servers?type=pitaya-admin", nil)
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, response.Code)
			},
		},
		"fail_no_server_of_type": {
			request: func() *http.Request {
				request, err := http.NewRequest("GET", "/servers?type=dontexist", nil)
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, response.Code)
				b := response.Body.String()
				assert.Equal(t, b, `{"success":false, "message":"failed to retrieve servers", "reason": "no servers available of this type"}`)
			},
		},
	}

	for name, table := range tables {
		t.Run(name, func(t *testing.T) {
			response = httptest.NewRecorder()
			serverHandler.ServeHTTP(response, table.request)
			table.asserts(response)
		})

	}
}
