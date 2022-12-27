package notifier

import (
	"context"
	"net/http"
	"net/url"

	"github.com/go-mail/mail"
	"github.com/krzysztof-gzocha/curnot/pkg/aggregator"
	"github.com/krzysztof-gzocha/curnot/pkg/config"
	"github.com/krzysztof-gzocha/curnot/pkg/formatter"
)

type Provider interface {
	GetNotifiers() []aggregator.Notifier
}

type Chain struct {
	client         *http.Client
	notifierConfig config.NotifiersConfig
}

func NewChain(c *http.Client, notifierConfig config.NotifiersConfig) *Chain {
	return &Chain{client: c, notifierConfig: notifierConfig}
}

func (c *Chain) Notify(ctx context.Context, msg aggregator.RateChange) error {
	notifiers, err := c.GetNotifiers(c.client)
	if err != nil {
		return err
	}

	for _, n := range notifiers {
		err := n.Notify(ctx, msg)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Chain) GetNotifiers(client *http.Client) ([]aggregator.Notifier, error) {
	var notifiers []aggregator.Notifier

	if emailNotifierConfig := c.notifierConfig.EmailConfig; emailNotifierConfig != nil {
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

	if c.notifierConfig.Desktop != nil {
		notifiers = append(notifiers, NewDesktop())
	}

	if httpConfig := c.notifierConfig.HttpConfig; httpConfig != nil {
		u, err := url.Parse(httpConfig.Path)
		if err != nil {
			return notifiers, err
		}

		notifiers = append(notifiers, NewHttp(
			client,
			*u,
			httpConfig.Method,
			httpConfig.AcceptedResponseStatuses,
			httpConfig.ExtraHeaders,
		))
	}

	return notifiers, nil
}
