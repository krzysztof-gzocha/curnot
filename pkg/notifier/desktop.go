// +build linux darwin freebsd

package notifier

import (
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

func (d *Desktop) Notify(msg string) error {
	return d.notifier.Push(
		aggregator.NotificationTitle,
		msg,
		"",
		notificator.UR_NORMAL,
	)
}
