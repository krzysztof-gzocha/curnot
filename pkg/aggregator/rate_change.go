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

	msg := "1 %s = %.3f %s"
	compare := r.New.Rate / r.Old.Rate
	if r.New.Rate > r.Old.Rate {
		msg = " (+" + msg
	} else if r.New.Rate <= r.Old.Rate {
		msg += " (-" + msg
	}

	return fmt.Sprintf(msg, compare, r.New.From, r.New.Rate, r.New.To)
}
