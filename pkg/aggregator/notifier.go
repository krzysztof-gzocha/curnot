package aggregator

import (
	"context"
)

const NotificationAppName = "Currency Notifier"
const NotificationTitle = "Currency alert"

type Notifier interface {
	Notify(ctx context.Context, msg RateChange) error
}
