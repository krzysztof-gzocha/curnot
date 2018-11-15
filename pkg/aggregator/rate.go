package aggregator

import (
	"fmt"
	"strings"

	"github.com/krzysztof-gzocha/curnot/pkg/config"
)

type Rate struct {
	From string
	To   string
	Rate float64
}

func (r *Rate) String() string {
	return fmt.Sprintf("%s%s", r.From, r.To)
}

func (r *Rate) supports(c config.CurrencyConfig) bool {
	return strings.EqualFold(r.From, c.From) && strings.EqualFold(r.To, c.To)
}

func (r *Rate) shouldNotify(alert config.Alert) bool {
	if alert.AnyChange {
		return true
	}

	if r.Rate >= alert.Above {
		return true
	}

	if r.Rate <= alert.Below {
		return true
	}

	return false
}
