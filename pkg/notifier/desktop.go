package notifier

import "github.com/0xAX/notificator"

type Desktop struct {
	notifier *notificator.Notificator
}

func NewDesktop(notifier *notificator.Notificator) *Desktop {
	return &Desktop{
		notifier: notifier,
	}
}

func (d *Desktop) Notify(msg string) error {
	return d.notifier.Push(
		"Currency notifier",
		msg,
		"",
		notificator.UR_NORMAL,
	)
}
