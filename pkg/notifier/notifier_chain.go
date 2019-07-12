package notifier

import (
	"github.com/go-mail/mail"
	"github.com/krzysztof-gzocha/curnot/pkg/aggregator"
	"github.com/krzysztof-gzocha/curnot/pkg/config"
	"github.com/krzysztof-gzocha/curnot/pkg/formatter"
)

type NotifierProvider interface {
	GetNotifiers() []aggregator.NotifierInterface
}

type NotifierChain struct {
	notifierConfig map[string]config.NotifierConfig
}

func NewNotifierChain(notifierConfig map[string]config.NotifierConfig) *NotifierChain {
	return &NotifierChain{notifierConfig: notifierConfig}
}

func (c *NotifierChain) Notify(msg string) error {
	notifiers := c.GetNotifiers()

	for _, n := range notifiers {
		err := n.Notify(msg)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *NotifierChain) GetNotifiers() []aggregator.NotifierInterface {
	var notifiers []aggregator.NotifierInterface

	emailNotifierConfig, exists := c.notifierConfig[NameEmailNotifier]

	if exists {
		dialer := mail.NewDialer(
			emailNotifierConfig.ConnectionParameters.Host,
			emailNotifierConfig.ConnectionParameters.Port,
			emailNotifierConfig.ConnectionParameters.Username,
			emailNotifierConfig.ConnectionParameters.Password,
		)

		notifier := NewEmail(
			*dialer,
			emailNotifierConfig.EmailReceiverParameters,
			emailNotifierConfig.ConnectionParameters,
			formatter.HtmlMessageFormatStrategy{},
		)

		notifiers = append(notifiers, notifier)
	}

	_, exists = c.notifierConfig[NameDesktopNotifier]

	if exists {
		notifiers = append(notifiers, NewDesktop())
	}

	return notifiers
}
