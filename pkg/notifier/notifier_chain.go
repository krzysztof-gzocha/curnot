package notifier

import (
	"context"
	"net/http"
	"net/url"

	"github.com/krzysztof-gzocha/curnot/pkg/aggregator"
	"github.com/krzysztof-gzocha/curnot/pkg/config"
	"github.com/wneessen/go-mail"
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
		client, err := mail.NewClient(
			emailNotifierConfig.ConnectionParameters.Host,
			mail.WithPort(emailNotifierConfig.ConnectionParameters.Port),
			mail.WithUsername(emailNotifierConfig.ConnectionParameters.Username),
			mail.WithPassword(emailNotifierConfig.ConnectionParameters.Password),
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithTLSPolicy(mail.TLSMandatory),
		)

		if err != nil {
			return notifiers, err
		}

		notifier := NewEmail(
			client,
			mail.TypeTextHTML,
			emailNotifierConfig.EmailReceiverParameters,
			emailNotifierConfig.ConnectionParameters,
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
