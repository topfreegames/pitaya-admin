package request

import (
	"github.com/topfreegames/pitaya/client"
)

// RequestMsg struct represents the message that the server will receive
type RequestMsg struct {
	Route     string
	Payload   []byte
	IsRequest bool
}

type clientMessageForwarder func(m interface{}) error

// ProcessMessage process received messages and send them to pitaya
func ProcessMessage(msg RequestMsg, pClient *client.Client) error {
	var err error
	if msg.IsRequest {
		_, err = pClient.SendRequest(msg.Route, msg.Payload)
	} else {
		err = pClient.SendNotify(msg.Route, msg.Payload)
	}

	return err
}

// ListenForClientMessages receives messages that were sent to the client and write it to the websocket
func ListenForClientMessages(pClient *client.Client, processMessage clientMessageForwarder) {
	for {
		select {
		case m := <-pClient.IncomingMsgChan:
			err := processMessage(m)
			if err != nil {
				return
			}
		}
	}
}
