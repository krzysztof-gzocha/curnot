package checker

import (
	"fmt"
	"time"
)

type TickerChecker struct {
	ticker  *time.Ticker
	checker Checker
}

func NewTickerChecker(
	ticker *time.Ticker,
	checker Checker,
) *TickerChecker {
	return &TickerChecker{
		ticker:  ticker,
		checker: checker,
	}
}

func (t *TickerChecker) Check() error {
	for range t.ticker.C {
		fmt.Println("Checking..")
		err := t.checker.Check()
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return nil
}
