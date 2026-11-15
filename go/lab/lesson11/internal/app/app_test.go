package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"testing"
	"time"
)

var discardLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

func listen(t *testing.T) net.Listener {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	return ln
}

func newClient() *http.Client {
	return &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
}

func TestNewServerHasTimeouts(t *testing.T) {
	srv := NewServer(http.NotFoundHandler())
	if srv.ReadHeaderTimeout == 0 || srv.ReadTimeout == 0 || srv.WriteTimeout == 0 || srv.IdleTimeout == 0 {
		t.Errorf("타임아웃이 모두 설정되어야 합니다: %+v", srv)
	}
}

func TestRunFinishesInFlightRequestBeforeReturning(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	slow := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		fmt.Fprint(w, "done")
	})

	ln := listen(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	runErr := make(chan error, 1)
	go func() {
		runErr <- Run(ctx, NewServer(slow), ln, 2*time.Second, discardLogger)
	}()

	type result struct {
		body string
		err  error
	}
	got := make(chan result, 1)
	go func() {
		resp, err := newClient().Get("http://" + ln.Addr().String())
		if err != nil {
			got <- result{err: err}
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		got <- result{body: string(body), err: err}
	}()

	<-started
	cancel()

	select {
	case err := <-runErr:
		t.Fatalf("요청이 끝나기 전에 Run이 돌아왔습니다: %v", err)
	case <-time.After(100 * time.Millisecond):
	}

	close(release)
	r := <-got
	if r.err != nil || r.body != "done" {
		t.Errorf("진행 중이던 요청이 정상 응답을 받아야 합니다: body=%q err=%v", r.body, r.err)
	}
	if err := <-runErr; err != nil {
		t.Errorf("Run이 nil을 돌려줘야 합니다: %v", err)
	}

	if resp, err := newClient().Get("http://" + ln.Addr().String()); err == nil {
		resp.Body.Close()
		t.Error("종료 뒤에는 새 연결이 거부되어야 합니다")
	}
}

func TestRunReturnsErrorWhenShutdownTimesOut(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	defer close(release)
	stuck := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
	})

	ln := listen(t)
	ctx, cancel := context.WithCancel(context.Background())
	runErr := make(chan error, 1)
	go func() {
		runErr <- Run(ctx, NewServer(stuck), ln, 50*time.Millisecond, discardLogger)
	}()

	go func() {
		resp, err := newClient().Get("http://" + ln.Addr().String())
		if err == nil {
			resp.Body.Close()
		}
	}()
	<-started
	cancel()

	err := <-runErr
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("got %v, want context.DeadlineExceeded", err)
	}
}

func TestRunReturnsServeError(t *testing.T) {
	ln := listen(t)
	ln.Close()

	err := Run(context.Background(), NewServer(http.NotFoundHandler()), ln, time.Second, discardLogger)
	if err == nil {
		t.Error("닫힌 리스너로는 서버를 시작할 수 없으므로 에러가 와야 합니다")
	}
}
