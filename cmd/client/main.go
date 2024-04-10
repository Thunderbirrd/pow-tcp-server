package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github/Thunderbirrd/pow-tcp-server/internal/client"
	"github/Thunderbirrd/pow-tcp-server/internal/config"
	"github/Thunderbirrd/pow-tcp-server/internal/pow"
)

func main() {
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stdout,
		zap.DebugLevel,
	))
	logger.Info("Client starting...")

	defer logger.Info("client stopped")

	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-ctx.Done()
		cancel()
	}()

	cfg := config.NewClientConfig()

	hashCash, err := pow.New(uint64(cfg.PowConfig.Complexity))
	if err != nil {
		logger.Error("failed to init pow", zap.Error(err))
		os.Exit(1)
	}

	c := client.New(cfg, logger, hashCash)
	if err = c.Start(ctx, cfg.MaxRequest); err != nil {
		logger.Error("failed to start client", zap.Error(err))
		os.Exit(1)
	}
}
