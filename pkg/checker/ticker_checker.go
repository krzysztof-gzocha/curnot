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

func (t *TickerChecker) Check(ctx context.Context) error {
	ticker := time.NewTicker(t.duration)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Checking..")
			timedCtx, cancel := context.WithTimeout(ctx, t.duration)
			err := t.checker.Check(timedCtx)
			cancel()
			if err != nil {
				fmt.Println(err.Error())
			}
		case <-ctx.Done():
			ticker.Stop()
			return ctx.Err()
		}
	}
}
