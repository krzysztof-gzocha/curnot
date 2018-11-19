// +build windows

package notifier

import (
	"github.com/go-toast/toast"
)

type Desktop struct{}

func NewDesktop() *Desktop {
	return &Desktop{}
}

func (d *Desktop) Notify(msg string) error {
	notification := toast.Notification{
		AppID:   NotificationAppName,
		Title:   NotificationTitle,
		Message: msg,
	}

	return notification.Push()
}
