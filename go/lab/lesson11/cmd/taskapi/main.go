package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golab/lesson11/internal/app"
	"golab/lesson11/internal/config"
	"golab/lesson11/internal/handler"
	"golab/lesson11/internal/repository"
	"golab/lesson11/internal/service"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "오류:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load(os.Getenv)
	if err != nil {
		return err
	}
	logger := newLogger(os.Stderr, cfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var repo service.Repository = repository.NewMemory()
	if cfg.DataFile != "" {
		fileRepo, err := repository.NewFile(cfg.DataFile)
		if err != nil {
			return err
		}
		repo = fileRepo
		logger.Info("파일 저장소를 사용합니다", "path", cfg.DataFile)
	}
	svc := service.New(repo)
	h := handler.New(svc, logger)
	root := handler.Chain(h.Routes(),
		handler.RequestID,
		handler.Logging(logger),
		handler.Recover(logger),
	)

	ln, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return fmt.Errorf("포트 %s를 열 수 없습니다: %w", cfg.Port, err)
	}
	return app.Run(ctx, app.NewServer(root), ln, cfg.ShutdownTimeout, logger)
}
