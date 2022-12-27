package checker

import (
	"context"
	"fmt"
	"time"
)

type TickerChecker struct {
	duration time.Duration
	checker  Checker
}

func NewTickerChecker(
	ticker time.Duration,
	checker Checker,
) *TickerChecker {
	return &TickerChecker{
		duration: ticker,
		checker:  checker,
	}
}

func (t *TickerChecker) StartChecking(ctx context.Context) {
	ticker := time.NewTicker(t.duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			timedCtx, cancel := context.WithTimeout(ctx, t.duration)
			err := t.checker.Check(timedCtx)
			cancel()
			if err != nil {
				fmt.Printf("=== Error: %+v\n", err)
			}
		case <-ctx.Done():
			fmt.Printf("=== Timed out after %s\n", t.duration)
		}
	}
}
