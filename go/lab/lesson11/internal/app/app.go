// Package app은 HTTP 서버를 실행하고 종료 신호가 오면 안전하게 멈추는 일을 맡습니다.
package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// NewServer는 운영에 쓸 만한 타임아웃을 설정한 http.Server를 만듭니다.
// http.ListenAndServe가 쓰는 기본 서버에는 이런 타임아웃이 없습니다.
func NewServer(h http.Handler) *http.Server {
	return &http.Server{
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

// Run은 ln에서 srv를 실행합니다. ctx가 취소되면 새 연결을 받지 않고,
// 진행 중인 요청이 끝나기를 shutdownTimeout까지 기다린 뒤 돌아옵니다.
func Run(ctx context.Context, srv *http.Server, ln net.Listener, shutdownTimeout time.Duration, logger *slog.Logger) error {
	serveErr := make(chan error, 1)
	go func() {
		logger.Info("서버 시작", "addr", ln.Addr().String())
		err := srv.Serve(ln)
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		serveErr <- err
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
	}

	logger.Info("종료 신호를 받았습니다. 진행 중인 요청을 마무리합니다", "timeout", shutdownTimeout)

	// ctx는 이미 취소되었으므로 그대로 쓰면 Shutdown이 기다리지 않고 바로 실패합니다.
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		_ = srv.Close()
		return fmt.Errorf("graceful shutdown 실패: %w", err)
	}
	if err := <-serveErr; err != nil {
		return err
	}
	logger.Info("서버가 정상 종료되었습니다")
	return nil
}
