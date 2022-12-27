//go:build linux || darwin || freebsd

package notifier

import (
	"context"

	"github.com/0xAX/notificator"
	"github.com/krzysztof-gzocha/curnot/pkg/aggregator"
)

type Desktop struct {
	notifier *notificator.Notificator
}

func NewDesktop() *Desktop {
	return &Desktop{
		notifier: notificator.New(notificator.Options{
			AppName: aggregator.NotificationAppName,
		}),
	}
}

func (d *Desktop) Notify(_ context.Context, msg aggregator.RateChange) error {
	return d.notifier.Push(
		aggregator.NotificationTitle,
		msg.String(),
		"",
		notificator.UR_NORMAL,
	)
}
