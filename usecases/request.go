package usecases

import (
	"github.com/gorilla/websocket"
	"github.com/topfreegames/pitaya/client"
)

// RequestMsg ...
type RequestMsg struct {
	Route     string
	Payload   []byte
	IsRequest bool
}

// ProcessMessage ...
func ProcessMessage(msg RequestMsg, pClient *client.Client) error {
	var err error
	if msg.IsRequest {
		_, err = pClient.SendRequest(msg.Route, msg.Payload)
	} else {
		err = pClient.SendNotify(msg.Route, msg.Payload)
	}
	if err != nil {
		return err
	}

	return nil
}

// ListenForClientMessages ...
func ListenForClientMessages(pClient *client.Client, ws *websocket.Conn) {
	for {
		select {
		case m := <-pClient.IncomingMsgChan:
			err := ws.WriteMessage(websocket.TextMessage, m.Data)
			if err != nil {
				return
			}
		}
	}
}
