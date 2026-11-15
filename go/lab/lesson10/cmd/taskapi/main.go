package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golab/lesson10/internal/app"
	"golab/lesson10/internal/config"
	"golab/lesson10/internal/handler"
	"golab/lesson10/internal/repository"
	"golab/lesson10/internal/service"
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

	repo := repository.NewMemory()
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
