package aggregator

import (
	"fmt"
)

type RateChange struct {
	Old *Rate
	New *Rate
}

func (r RateChange) String() string {
	if r.Old == nil {
		return fmt.Sprintf(
			"1 %s = %.3f %s",
			r.New.From,
			r.New.Rate,
			r.New.To,
		)
	}

	msg := ""
	compare := r.New.Rate / r.Old.Rate
	if r.New.Rate > r.Old.Rate {
		msg += " (+"
	} else {
		msg += " (-"
	}
	msg += "%.2f) 1 %s = %.3f %s"

	return fmt.Sprintf(msg, compare, r.New.From, r.New.Rate, r.New.To)
}
