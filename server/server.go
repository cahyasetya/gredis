package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/cahyasetya/gredis/constants"
	"github.com/cahyasetya/gredis/logger"
	"github.com/cahyasetya/gredis/parser"
	"github.com/cahyasetya/gredis/processors"
	"github.com/cahyasetya/gredis/types"
	"go.uber.org/zap"
)

type Config struct {
	Port int
}

type Server struct {
	config Config
	listener net.Listener
	wg sync.WaitGroup
}

func New(config Config) *Server {
	return &Server{
		config: config,
	}
}

func (s *Server) Start(ctx context.Context) error {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	logger.Log.Info("Server started", zap.Int("port", s.config.Port))

	go s.acceptConnections(ctx)

	return nil
}

func (s *Server) acceptConnections(ctx context.Context) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					logger.Log.Warn("Temporary error when accepting connection", zap.Error(err))
					time.Sleep(time.Millisecond * 100) // Short delay before retrying
					continue
				}
				logger.Log.Error("Failed to accept connection", zap.Error(err))
				return // Exit the loop if it's a non-temporary error
			}
		}

		fmt.Println("Got new connection")
		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	buffer := make([]byte, constants.BufferSize)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				logger.Log.Info("Client disconnected")
			} else {
				logger.Log.Error("Error reading from connection", zap.Error(err))
			}
			return
		}

		message := types.Message(buffer[:n])
		commands := parser.SplitCommand(message)

		logger.Log.Info("Received commands",
			zap.Int("count", len(commands)),
			zap.Strings("commands", commands))

		result := processors.HandleCommand(commands)
		_, err = conn.Write(result)
		if err != nil {
			logger.Log.Error("Error writing response", zap.Error(err))
		}
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.listener != nil {
		s.listener.Close()
	}

	// Signal all goroutines to stop
	// (You might need to implement this signaling mechanism)
	// s.stopAllGoroutines()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
