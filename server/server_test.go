package server

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	config := Config{Port: 6379}
	srv := New(config)

	assert.NotNil(t, srv)
	assert.Equal(t, config, srv.config)
}

func TestServerStart(t *testing.T) {
	config := Config{Port: 0} // Use port 0 to let the system assign a free port
	srv := New(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := srv.Start(ctx)
	require.NoError(t, err)

	// Check if the server is listening
	assert.NotNil(t, srv.listener)
	assert.NotEqual(t, 0, srv.listener.Addr().(*net.TCPAddr).Port)

	// Try to connect to the server
	conn, err := net.Dial("tcp", srv.listener.Addr().String())
	require.NoError(t, err)
	conn.Close()
}

func TestServerShutdown(t *testing.T) {
	config := Config{Port: 0}
	srv := New(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := srv.Start(ctx)
	require.NoError(t, err)

	// Start a long-running connection
	conn, err := net.Dial("tcp", srv.listener.Addr().String())
	require.NoError(t, err)
	defer conn.Close()

	// Shutdown the server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	err = srv.Shutdown(shutdownCtx)
	assert.NoError(t, err)

	// Try to connect after shutdown
	_, err = net.Dial("tcp", srv.listener.Addr().String())
	assert.Error(t, err)
}

func TestHandleConnection(t *testing.T) {
	config := Config{Port: 0}
	srv := New(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := srv.Start(ctx)
	require.NoError(t, err)

	conn, err := net.Dial("tcp", srv.listener.Addr().String())
	require.NoError(t, err)
	defer conn.Close()

	// Send a test command
	_, err = conn.Write([]byte("PING\r\n"))
	require.NoError(t, err)

	// Read the response
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	require.NoError(t, err)

	// Check the response
	assert.Equal(t, "+PONG\r\n", string(buffer[:n]))
}
