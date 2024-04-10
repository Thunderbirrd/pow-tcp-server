package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"

	"github/Thunderbirrd/pow-tcp-server/internal/config"
	"github/Thunderbirrd/pow-tcp-server/pkg/utils"
)

const protocol = "tcp"

type Server struct {
	cfg      *config.ServerConfig
	logger   *zap.Logger
	verifier pow
	repo     repo
	listener net.Listener
	wg       sync.WaitGroup
	cancel   context.CancelFunc
}

func New(cfg *config.ServerConfig, logger *zap.Logger, verifier verifier, repo repo) *Server {
	return &Server{
		cfg:      cfg,
		logger:   logger,
		verifier: verifier,
		repo:     repo,
	}
}

func (s *Server) Run(ctx context.Context) (err error) {
	ctx, s.cancel = context.WithCancel(ctx)
	defer s.cancel()

	listenConfig := net.ListenConfig{
		KeepAlive: s.cfg.KeepAlive,
	}

	s.listener, err = listenConfig.Listen(ctx, protocol, s.cfg.Address)
	if err != nil {
		s.logger.Error("failed to listen", zap.Error(err))
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.logger.Info(fmt.Sprintf("Server started with address: %s", s.cfg.Address))

	s.wg.Add(1)
	go s.serve(ctx)
	s.wg.Wait()

	s.logger.Info("server stopped")

	return nil
}

// Stop stops the server
func (s *Server) Stop() {
	s.cancel()
}

func (s *Server) serve(ctx context.Context) {
	defer s.wg.Done()
	go func() {
		<-ctx.Done()
		err := s.listener.Close()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			s.logger.Error("failed to close listener: ", zap.Error(err))
		}
	}()

	for {
		conn, err := s.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			s.logger.Debug("listener closed")
			return
		} else if err != nil {
			s.logger.Error("failed to accept connection: ", zap.Error(err))
			continue
		}

		s.wg.Add(1)
		go func(conn net.Conn) {
			defer s.wg.Done()

			if err = s.handle(conn); err != nil {
				s.logger.Error("internal error: ", zap.Error(err))
			}
		}(conn)
	}
}

func (s *Server) handle(conn net.Conn) error {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			s.logger.Error("failed to close connection: ", zap.Error(err))
		}
	}(conn)

	err := conn.SetDeadline(time.Now().Add(s.cfg.Deadline))
	if err != nil {
		s.logger.Error("failed to set conn deadline: ", zap.Error(err))
	}

	if _, err = utils.ReadMessage(conn); err != nil {
		return fmt.Errorf("read message err: %w", err)
	}

	challenge := s.verifier.Challenge()
	if err := utils.WriteMessage(conn, challenge); err != nil {
		return fmt.Errorf("challenge err: %w", err)
	}

	solution, err := utils.ReadMessage(conn)
	if err != nil {
		return fmt.Errorf("receive pow err: %w", err)
	}

	if err = s.verifier.Verify(challenge, solution); err != nil {
		return fmt.Errorf("solution verification err: %w", err)
	}

	line, err := s.repo.GetLine()
	if err != nil {
		return fmt.Errorf("get line from repo err: %w", err)
	}

	if err = utils.WriteMessage(conn, []byte(line)); err != nil {
		return fmt.Errorf("send line err: %w", err)
	}

	return nil
}
