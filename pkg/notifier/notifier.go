package notifier

const NotificationAppName = "Currency Notifier"
const NotificationTitle = "Currency alert"

type NotifierInterface interface {
	Notify(msg string) error
}
