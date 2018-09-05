package api_test

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/pitaya-admin/api"
)

func TestNewApp(t *testing.T) {
	t.Parallel()
	App, err := api.NewApp("0.0.0.0", 8000, viper.New())
	assert.NoError(t, err)
	assert.NotNil(t, App)
}
