//go:build windows

package notifier

import (
	"context"

	"github.com/go-toast/toast"
	"github.com/krzysztof-gzocha/curnot/pkg/aggregator"
)

type Desktop struct{}

func NewDesktop() *Desktop {
	return &Desktop{}
}

func (d *Desktop) Notify(ctx context.Context, msg aggregator.RateChange) error {
	notification := toast.Notification{
		// https://github.com/go-toast/toast/issues/9
		AppID:   "{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\\WindowsPowerShell\\v1.0\\powershell.exe",
		Title:   aggregator.NotificationTitle,
		Message: msg.String(),
	}

	return notification.Push()
}
