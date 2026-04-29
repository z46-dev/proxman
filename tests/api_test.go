package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/z46-dev/proxman"
)

func TestAPIFailsWithIncorrectStuff(t *testing.T) {
	mustConfig(t)

	var (
		client *proxman.Client
		err    error
	)

	// 127.0.0.0/8 is reserved for loopback, so this should fail to connect to any API
	client, err = proxman.NewClient("https://127.127.127.127:8006", "token", "token")

	assert.Error(t, err)
	assert.NotNil(t, client)
	assert.False(t, client.Living())
}

func TestAPILoadNodes(t *testing.T) {
	mustConfig(t)

	if !config.EnableProxmoxTesting {
		t.Skip("Proxmox API testing is disabled in config")
	}
}
