package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/minhaz11/cache"
	"go.uber.org/zap"
)

const (
	maxWorkers = 100
	port       = ":6369"
)

type Server struct {
	cache    *cache.Cache
	Logger   *zap.SugaredLogger
	listener net.Listener
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewServer() (*Server, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}
	sugar := logger.Sugar()

	c, err := cache.NewCache("cache.json")
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		cache:  c,
		Logger: sugar,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (s *Server) startListener() error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}
	s.listener = ln
	s.Logger.Infof("Cache server is listening on port %s...", port)
	return nil
}

func (s *Server) setupWorkerPool(connectionChan chan net.Conn) *sync.WaitGroup {
	var wg sync.WaitGroup

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-s.ctx.Done():
					return
				case conn, ok := <-connectionChan:
					if !ok {
						return
					}

					_ = conn.SetReadDeadline(time.Now().Add(30 * time.Second))
					_ = conn.SetWriteDeadline(time.Now().Add(30 * time.Second))

					s.cache.HandleConnection(conn)
				}
			}
		}()
	}

	return &wg
}

func (s *Server) setupShutdownHandler(connectionChan chan net.Conn, wg *sync.WaitGroup) {
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		s.Logger.Info("Shutting down server...")
		s.cancel()
		s.listener.Close()
		close(connectionChan)
		wg.Wait()
	}()
}

func (s *Server) acceptConnections(connectionChan chan net.Conn) error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.ctx.Done():
				return nil
			default:
				s.Logger.Errorf("Error accepting connection: %v", err)
				continue
			}
		}
		connectionChan <- conn
	}
}

func (s *Server) Run() error {
	connectionChan := make(chan net.Conn, maxWorkers)

	if err := s.startListener(); err != nil {
		return err
	}
	defer s.listener.Close()

	wg := s.setupWorkerPool(connectionChan)

	s.setupShutdownHandler(connectionChan, wg)

	return s.acceptConnections(connectionChan)
}
