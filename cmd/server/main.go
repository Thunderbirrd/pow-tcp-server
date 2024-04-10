package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github/Thunderbirrd/pow-tcp-server/internal/config"
	"github/Thunderbirrd/pow-tcp-server/internal/pow"
	"github/Thunderbirrd/pow-tcp-server/internal/repository"
	"github/Thunderbirrd/pow-tcp-server/internal/server"
)

func main() {
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stdout,
		zap.DebugLevel,
	))
	logger.Info("Server starting...")

	defer logger.Info("server stopped")

	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-ctx.Done()
		cancel()
	}()

	cfg := config.NewServerConfig()

	hashCash, err := pow.New(uint64(cfg.PowConfig.Complexity))
	if err != nil {
		logger.Error("failed to init pow", zap.Error(err))
		os.Exit(1)
	}

	repo := repository.New()
	srv := server.New(cfg, logger, hashCash, repo)

	if err = srv.Run(ctx); err != nil {
		logger.Error("server error", zap.Error(err))
		os.Exit(1)
	}
}
