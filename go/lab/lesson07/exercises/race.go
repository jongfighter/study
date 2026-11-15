package exercises

import (
	"context"
	"time"
)

// Source는 결과 문자열 하나를 만드는 작업입니다. ctx가 취소되면 서둘러 돌아와야 합니다.
type Source func(ctx context.Context) string

// Race는 sources를 동시에 실행해 가장 먼저 끝난 결과를 돌려줍니다.
// timeout 안에 아무도 끝나지 않으면 ctx.Err()를 돌려주고, 늦은 작업은 취소됩니다.
func Race(ctx context.Context, timeout time.Duration, sources ...Source) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	results := make(chan string, len(sources))
	for _, src := range sources {
		go func() {
			results <- src(ctx)
		}()
	}

	select {
	case r := <-results:
		return r, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
