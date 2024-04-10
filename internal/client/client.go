package client

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net"

	"github/Thunderbirrd/pow-tcp-server/internal/config"
	"github/Thunderbirrd/pow-tcp-server/pkg/utils"
)

const protocol = "tcp"

type Client struct {
	conf   *config.ClientConfig
	logger *zap.Logger
	solver pow
}

func New(conf *config.ClientConfig, logger *zap.Logger, solver pow) *Client {
	return &Client{
		conf:   conf,
		logger: logger,
		solver: solver,
	}
}

func (c *Client) Start(ctx context.Context, count int) error {
	for i := 0; i < count; i++ {
		if ctx.Err() != nil {
			break
		}

		q, err := c.GetLine(ctx)
		if err != nil {
			c.logger.Error("failed to get message: ", zap.Error(err))
		} else {
			c.logger.Info(string(q))
		}
	}

	return nil
}

func (c *Client) GetLine(ctx context.Context) ([]byte, error) {
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, protocol, c.conf.ServerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %w", err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			c.logger.Error("failed to close connection: ", zap.Error(err))
		}
	}(conn)

	if err := utils.WriteMessage(conn, []byte("challenge")); err != nil {
		return nil, fmt.Errorf("send request err: %w", err)
	}

	challenge, err := utils.ReadMessage(conn)
	if err != nil {
		return nil, fmt.Errorf("receive challenge err: %w", err)
	}

	solution := c.solver.Solve(challenge)
	if err := utils.WriteMessage(conn, solution); err != nil {
		return nil, fmt.Errorf("send solution err: %w", err)
	}

	line, err := utils.ReadMessage(conn)
	if err != nil {
		return nil, fmt.Errorf("receive line err: %w", err)
	}

	return line, nil
}
