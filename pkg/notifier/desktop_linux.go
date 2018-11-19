// +build linux darwin freebsd

package notifier

import (
	"github.com/0xAX/notificator"
)

type Desktop struct {
	notifier *notificator.Notificator
}

func NewDesktop() *Desktop {
	return &Desktop{
		notifier: notificator.New(notificator.Options{
			AppName: NotificationAppName,
		}),
	}
}

func (d *Desktop) Notify(msg string) error {
	return d.notifier.Push(
		NotificationTitle,
		msg,
		"",
		notificator.UR_NORMAL,
	)
}
