package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/topfreegames/pitaya-admin/request"
	"github.com/topfreegames/pitaya/client"
	"github.com/topfreegames/pitaya/logger"
)

// RequestHandler handle requests
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
		WriteError(w, http.StatusInternalServerError, "failed to establish websocket", err)
		return
	}

	defer func() {
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		ws.Close()
	}()

	ws.SetReadDeadline(time.Now().Add(s.readDeadLine))

	var pClient = client.New(logrus.InfoLevel)
	err = pClient.ConnectTo(pitayaAddress)

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to connect to pitaya", err)
		ws.Close()
		return
	}

	defer pClient.Disconnect()

	go request.ListenForClientMessages(pClient, func(m interface{}) error {
		byteMessage, err := json.Marshal(m)
		if err != nil {
			return err
		}
		ws.WriteMessage(websocket.BinaryMessage, byteMessage)
		return nil
	})

	for {
		select {
		default:
			var msg request.RequestMsg
			if err := ws.ReadJSON(&msg); err != nil {
				logger.Log.Errorf("failed to read msg from websocket, closing connection. %s", err.Error())
				return
			}
			request.ProcessMessage(msg, pClient)
		}

	}
}
