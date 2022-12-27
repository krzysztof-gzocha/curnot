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

const NameNotifierDesktop = "desktop"
const NameNotifierEmail = "email"
const NameNotifierHttp = "http"

type Provider interface {
	GetNotifiers() []aggregator.Notifier
}

type Chain struct {
	client         *http.Client
	notifierConfig map[string]config.NotifierConfig
}

func NewChain(c *http.Client, notifierConfig map[string]config.NotifierConfig) *Chain {
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

	emailNotifierConfig, exists := c.notifierConfig[NameNotifierEmail]

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

	_, exists = c.notifierConfig[NameNotifierDesktop]

	if exists {
		notifiers = append(notifiers, NewDesktop())
	}

	httpConfig, exists := c.notifierConfig[NameNotifierHttp]
	if exists {
		u, err := url.Parse(httpConfig.HttpParameters.Path)
		if err != nil {
			return notifiers, err
		}

		notifiers = append(notifiers, NewHttp(
			client,
			*u,
			httpConfig.HttpParameters.Method,
			httpConfig.HttpParameters.AcceptedResponseStatuses,
		))
	}

	return notifiers, nil
}
