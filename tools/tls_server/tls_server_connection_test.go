package main

import (
	"../../client_server_connection/core/handle_connections"
	"../tls_server"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTlsServerListening(t *testing.T) {
	localip, err := handle_connections.GetLocalIp04Str()
	assert.Nil(t, err)
	localip = localip + ":9090"
	main.StartTLSServerTo(localip)
}
