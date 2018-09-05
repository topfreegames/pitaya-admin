package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/pitaya-admin/api"
	pt "github.com/topfreegames/pitaya-admin/testing"
)

func TestPushHandler(t *testing.T) {
	t.Parallel()
	if !pt.IsConf {
		pt.ConfApp()
	}
	a := pt.TestApp
	var (
		response    *httptest.ResponseRecorder
		pushHandler = api.NewPushToUsersHandler(a)
	)

	tables := map[string]struct {
		request *http.Request
		asserts func(response *httptest.ResponseRecorder)
	}{
		"test_empty_body": {
			request: func() *http.Request {
				request, err := http.NewRequest("POST", "/user/push", nil)
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
				body := "not a pushmsg struct"
				bts, _ := json.Marshal(body)

				request, err := http.NewRequest("POST", "/user/push", bytes.NewReader(bts))
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, response.Code)
				b := response.Body.String()
				assert.Equal(t, `{"success":false, "message":"failed to decode request body into push struct", "reason": "json: cannot unmarshal string into Go value of type api.PushMsg"}`, b)

			},
		},
		"test_bad_route": {
			request: func() *http.Request {
				body := &api.PushMsg{
					Uids:         []string{"1", "2"},
					Route:        "ro.ut.e",
					Message:      1,
					FrontendType: "",
				}
				bts, _ := json.Marshal(body)

				request, err := http.NewRequest("POST", "/user/push", bytes.NewReader(bts))
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, response.Code)
			},
		},
		"test_success": {
			request: func() *http.Request {
				body := &api.PushMsg{
					Uids:         []string{"1", "2"},
					Route:        "ro.u.te",
					Message:      2,
					FrontendType: "connector",
				}
				bts, err := json.Marshal(body)
				assert.NoError(t, err)
				request, err := http.NewRequest("POST", "/user/push", bytes.NewReader(bts))
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, response.Code)
			},
		},
		// TODO:   make server type and test
	}

	for name, table := range tables {
		t.Run(name, func(t *testing.T) {
			response = httptest.NewRecorder()
			pushHandler.ServeHTTP(response, table.request)
			table.asserts(response)
		})
	}

}

func TestKickHandler(t *testing.T) {
	// t.Parallel()
	if !pt.IsConf {
		pt.ConfApp()
	}
	a := pt.TestApp
	var (
		response    *httptest.ResponseRecorder
		kickHandler = api.NewKickUserHandler(a)
	)

	tables := map[string]struct {
		request *http.Request
		asserts func(response *httptest.ResponseRecorder)
	}{
		"test_empty_body": {
			request: func() *http.Request {
				request, err := http.NewRequest("POST", "/user/kick", nil)
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, response.Code)
				b := response.Body.String()
				assert.Equal(t, `{"success":false, "message":"request body shouldnt be empty", "reason": "empty request body"}`, b)

			},
		},
		"test_bad_post": {
			request: func() *http.Request {
				body := "not a kick msg struct"
				bts, _ := json.Marshal(body)

				request, err := http.NewRequest("POST", "/user/kick", bytes.NewReader(bts))
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				b := response.Body.String()
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Equal(t, `{"success":false, "message":"failed to decode request body into kick struct", "reason": "json: cannot unmarshal string into Go value of type api.KickMsg"}`, b)
			},
		},
		"test_bad_route": {
			request: func() *http.Request {
				body := &api.KickMsg{
					Uids:         []string{"1", "2"},
					FrontendType: "",
				}
				bts, err := json.Marshal(body)
				assert.NoError(t, err)
				request, err := http.NewRequest("POST", "/user/kick", bytes.NewReader(bts))
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, response.Code)
			},
		},
		"test_success": {
			request: func() *http.Request {
				body := &api.KickMsg{
					Uids:         []string{"1", "2"},
					FrontendType: "connector",
				}
				bts, err := json.Marshal(body)
				assert.NoError(t, err)
				request, err := http.NewRequest("POST", "/user/kick", bytes.NewReader(bts))
				assert.NoError(t, err)
				return request
			}(),
			asserts: func(response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, response.Code)
			},
		},
	}

	for name, table := range tables {
		t.Run(name, func(t *testing.T) {
			response = httptest.NewRecorder()
			kickHandler.ServeHTTP(response, table.request)
			table.asserts(response)
		})
	}
}
