package notifier

type NotifierInterface interface {
	Notify(msg string) error
}
