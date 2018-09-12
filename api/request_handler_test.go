package api_test

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/pitaya-admin/api"
	"github.com/topfreegames/pitaya-admin/request"
	pt "github.com/topfreegames/pitaya-admin/testing"
)

// Message struct as sent by pitaya
type Message struct {
	Data []byte
}

func TestRequestHandler(t *testing.T) {
	t.Parallel()

	if !pt.IsConf {
		pt.ConfApp()
	}
	a := pt.TestApp

	s := httptest.NewServer(api.NewRequestHandler(a))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http") + "/request?address=localhost:3250"

	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.NoError(t, err)
	defer ws.Close()

	firstRequest := &request.RequestMsg{
		Route:     "connector.setsessiondata",
		Payload:   []byte(`{"data":{"key1":"value1"}}`),
		IsRequest: false,
	}
	err = ws.WriteJSON(firstRequest)
	assert.NoError(t, err)

	secondRequest := &request.RequestMsg{
		Route:     "connector.getsessiondata",
		Payload:   nil,
		IsRequest: true,
	}
	err = ws.WriteJSON(secondRequest)
	assert.NoError(t, err)

	_, p, err := ws.ReadMessage()
	assert.NoError(t, err)
	pmsg := &Message{}

	err = json.Unmarshal(p, pmsg)
	assert.NotNil(t, pmsg)
	assert.NoError(t, err)

}
