package docs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/pitaya-admin/docs"
	pt "github.com/topfreegames/pitaya-admin/testing"
)

func TestGetDocumentationForServerType(t *testing.T) {
	if !pt.IsConf {
		pt.ConfApp()
	}
	tables := []struct {
		name            string
		serverType      string
		docsRemoteRoute string
		err             bool
	}{
		{"success", "connector", "remote.docs", false},
		{"fail_no_route_in_server", "room", "remote.docs", true},
		{"fail_wrong_route", "connector", "doc.docs", true},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			ret, err := docs.GetDocumentationForServerType(table.serverType, table.docsRemoteRoute, true)
			if table.err {
				assert.Nil(t, ret)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, ret)
				assert.NoError(t, err)
			}
		})

	}
}

func TestGetDocForHandlersOrRemotes(t *testing.T) {

	if !pt.IsConf {
		pt.ConfApp()
	}

	tables := []struct {
		name            string
		serverType      string
		docsRemoteRoute string
		handlerOrRemote string
	}{
		// Fail scenarios are tested in above function
		{"succes_handler", "connector", "remote.docs", "handler"},
		{"succes_remote", "connector", "remote.docs", "remote"},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			ret, err := docs.GetDocForHandlersOrRemotes(table.serverType, table.handlerOrRemote, table.docsRemoteRoute, true)
			assert.NoError(t, err)
			assert.NotNil(t, ret)
		})
	}
}

func TestGetMethodDoc(t *testing.T) {

	if !pt.IsConf {
		pt.ConfApp()
	}

	tables := []struct {
		name            string
		serverType      string
		methodType      string
		route           string
		docsRemoteRoute string
	}{
		{"success", "connector", "remote", "connector.connectorremote.docs", "remote.docs"},
		{"fail", "connector", "remote", "connector.getsessiondata", "remote.docs"},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			ret, err := docs.GetMethodDoc(table.serverType, table.methodType, table.route, table.docsRemoteRoute, true)
			assert.NoError(t, err)
			if table.name == "fail" {
				//Method not found will return a nil map
				assert.Nil(t, ret)
			} else {
				assert.NotNil(t, ret)
			}
		})
	}
}
