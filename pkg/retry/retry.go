package retry

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

func Retry(attempts int, sleep time.Duration, f func() error) (err error) {
	return WithCancel(context.Background(), attempts, sleep, f)
}

func WithCancel(ctx context.Context, attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; i < attempts; i++ {
		if err = f(); err == nil {
			// success
			return
		}

		// retry after duration X
		sleepWithJitter := wait.Jitter(sleep, 0.6)
		fmt.Printf("retrying in %s after error: %v\n", sleepWithJitter, err)

		if ctx.Err() != nil {
			return ctx.Err()
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(sleepWithJitter):
			continue
		}
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
