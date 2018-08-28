package api

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/topfreegames/pitaya-admin/usecases"
	"github.com/topfreegames/pitaya/client"
	"github.com/topfreegames/pitaya/logger"
)

// RequestMsg struct
type RequestMsg struct {
	Route     string
	Payload   []byte
	IsRequest bool
}

// RequestHandler handler
type RequestHandler struct {
	App          *App
	readDeadLine time.Duration
}

// NewRequestHandler creates a new request handler
func NewRequestHandler(a *App) *RequestHandler {
	m := &RequestHandler{
		App:          a,
		readDeadLine: a.Config.GetDuration("pitayaadmin.request.readdeadline"),
	}
	return m
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ServeHTTP creates a web socket connection with the user
func (s *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	pitayaAddress := v.Get("address")
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		WriteError(w, 500, "failed to establish websocket", err)
		return
	}

	defer func() {
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		ws.Close()
	}()

	ws.SetReadDeadline(time.Now().Add(s.readDeadLine))

	var pClient = client.New(logrus.InfoLevel)
	err = pClient.ConnectTo(pitayaAddress)
	defer pClient.Disconnect()

	if err != nil {
		WriteError(w, 500, "failed to connect to pitaya", err)
		ws.Close()
		return
	}

	// TODO: this should receive a function instead of the websocket itself
	go usecases.ListenForClientMessages(pClient, ws)

	for { //read requests and notify`s
		select {
		default:
			var msg usecases.RequestMsg
			if err := ws.ReadJSON(&msg); err != nil {
				logger.Log.Warnf("failed to read msg from websocket, closing connection. %s", err.Error())
				return
			}
			usecases.ProcessMessage(msg, pClient)
		}

	}
}
