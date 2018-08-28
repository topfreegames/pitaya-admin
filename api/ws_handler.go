package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/topfreegames/pitaya/client"
	"github.com/topfreegames/pitaya/logger"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// RequestMsg struct
type RequestMsg struct {
	Route     string
	Payload   []byte
	IsRequest bool
}

// RequestHandler handler
type RequestHandler struct {
	App *App
}

// NewRequestHandler creates a new request handler
func NewRequestHandler(a *App) *RequestHandler {
	m := &RequestHandler{App: a}
	return m
}

// ServeHTTP creates a web socket connection with the user
func (s *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		WriteError(w, 500, "failed to establish websocket", err)
		return
	}
	logger.Log.Info("Successfully upgraded connection to websocket")

	defer func() {
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		ws.Close()
	}()

	ws.SetReadDeadline(time.Now().Add(15 * time.Second)) // this shouldnt be hard coded
	logger.Log.Info("Connecting to pitaya")

	var pClient = client.New(logrus.InfoLevel)
	err = pClient.ConnectTo("localhost:3250")
	defer pClient.Disconnect()

	if err != nil {
		WriteError(w, 500, "failed to connect to pitaya", err)
		ws.Close()
		return
	}
	logger.Log.Info("Successfully connected to pitaya")

	go func() { //forward pitaya client
		for {
			select {
			case m := <-pClient.IncomingMsgChan:
				err := ws.WriteMessage(websocket.TextMessage, m.Data)
				if err != nil {
					return
				}
			}
		}
	}()

	for { //read requests and notify`s
		select {
		default:
			var msg RequestMsg
			if err := ws.ReadJSON(&msg); err != nil {
				logger.Log.Warnf("failed to read msg from websocket, closing connection. %s", err.Error())
				return
			}
			if msg.IsRequest {
				_, err = pClient.SendRequest(msg.Route, msg.Payload)
			} else {
				err = pClient.SendNotify(msg.Route, msg.Payload)
			}
			if err != nil {
				fmt.Println("error.: ", err)
				return
			}
		}

	}
}
