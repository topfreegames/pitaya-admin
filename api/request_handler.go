package api

import (
	"encoding/json"
	"net/http"
	"strings"
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
	whiteList    []string
}

// NewRequestHandler creates a new request handler
func NewRequestHandler(a *App) *RequestHandler {
	m := &RequestHandler{
		App:          a,
		readDeadLine: a.Config.GetDuration("request.readdeadline"),
		whiteList:    a.Config.GetStringSlice("request.whitelist"),
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
	upgrader.CheckOrigin = func(r *http.Request) bool {
		if strings.Contains(r.Host, "127.0.0.1:") {
			return true
		}
		for _, a := range s.whiteList {
			if a == r.Host {
				return true
			}
		}
		return false
	}
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to establish websocket", err)
		return
	}

	defer func() {
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		logger.Log.Info("Client disconnected")
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

	logger.Log.Info("connected to pitaya")

	defer pClient.Disconnect()

	go request.ListenForClientMessages(pClient, func(m interface{}) error {
		byteMessage, err := json.Marshal(m)
		if err != nil {
			return err
		}
		ws.WriteMessage(websocket.TextMessage, byteMessage)
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
