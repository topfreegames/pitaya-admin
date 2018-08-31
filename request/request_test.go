package request_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/pitaya-admin/request"
	pt "github.com/topfreegames/pitaya-admin/testing"
	"github.com/topfreegames/pitaya/client"
)

type Message struct {
	Type       byte   // message type
	ID         uint   // unique id, zero while notify mode
	Route      string // route for locating service
	Data       []byte // payload
	compressed bool   // is message compressed
	Err        bool   // is an error message
}

func TestProcessMessage(t *testing.T) {

	if pt.IsConf == false {
		pt.ConfApp()
	}

	m1 := request.RequestMsg{
		Route:     "connector.connectorremote.remotefunc",
		Payload:   []byte("bla"),
		IsRequest: true,
	}

	var pClient = client.New(logrus.InfoLevel)
	err := pClient.ConnectTo("localhost:3250")
	assert.NoError(t, err)
	defer pClient.Disconnect()

	err = request.ProcessMessage(m1, pClient)
	assert.NoError(t, err)

	m2 := request.RequestMsg{
		Route:     "connector.connectorremote.remotefunc",
		Payload:   []byte("bla"),
		IsRequest: false,
	}

	err = request.ProcessMessage(m2, pClient)
	assert.NoError(t, err)

}

func TestListenForClientMessages(t *testing.T) {
	var pClient = client.New(logrus.InfoLevel)
	err := pClient.ConnectTo("localhost:3250")
	assert.NoError(t, err)
	defer pClient.Disconnect()

	go request.ListenForClientMessages(pClient, func(m interface{}) error {
		byteMessage, err := json.Marshal(m)
		assert.NoError(t, err)
		fmt.Println(byteMessage)
		return nil
	})

}
