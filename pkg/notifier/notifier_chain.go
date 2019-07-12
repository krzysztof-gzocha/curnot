package notifier

import (
	"github.com/go-mail/mail"
	"github.com/krzysztof-gzocha/curnot/pkg/aggregator"
	"github.com/krzysztof-gzocha/curnot/pkg/config"
	"github.com/krzysztof-gzocha/curnot/pkg/formatter"
)

const NameDesktopNotifier = "desktop"
const NameEmailNotifier = "email"

type Provider interface {
	GetNotifiers() []aggregator.NotifierInterface
}

type Chain struct {
	notifierConfig map[string]config.NotifierConfig
}

func NewChain(notifierConfig map[string]config.NotifierConfig) *Chain {
	return &Chain{notifierConfig: notifierConfig}
}

func (c *Chain) Notify(msg string) error {
	notifiers := c.GetNotifiers()

	for _, n := range notifiers {
		err := n.Notify(msg)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Chain) GetNotifiers() []aggregator.NotifierInterface {
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
