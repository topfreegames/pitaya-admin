package request_test

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/pitaya-admin/request"
	pt "github.com/topfreegames/pitaya-admin/testing"
	"github.com/topfreegames/pitaya/client"
)

func TestProcessMessage(t *testing.T) {

	if !pt.IsConf {
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
	var wg sync.WaitGroup
	wg.Add(1)
	go request.ListenForClientMessages(pClient, func(m interface{}) error {

		msg, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.NotNil(t, msg)
		wg.Done()
		return nil
	})
	m2 := request.RequestMsg{
		Route:     "connector.connectorremote.getsessiondata",
		Payload:   nil,
		IsRequest: true,
	}
	err = request.ProcessMessage(m2, pClient)
	assert.NoError(t, err)
	wg.Wait()
}
