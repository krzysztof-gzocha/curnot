package notifier

type NotifierChain struct {
	notifierProvider NotifierProvider
}

func NewNotifierChain(notifierProvider *NotifierProvider) *NotifierChain {
	return &NotifierChain{notifierProvider: *notifierProvider}
}

func (c *NotifierChain) Notify(msg string) error {
	notifiers := c.notifierProvider.GetNotifiers()

	for _, n := range notifiers {
		err := n.Notify(msg)

		if err != nil {
			return err
		}
	}

	return nil
}
