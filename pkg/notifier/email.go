package notifier

import (
	"context"
	"html/template"
	"log"

	"github.com/krzysztof-gzocha/curnot/pkg/aggregator"
	"github.com/krzysztof-gzocha/curnot/pkg/config"
	"github.com/wneessen/go-mail"
)

const bodyTemplate = "<p>{{ . }}</p>"

type Email struct {
	client               *mail.Client
	contentType          mail.ContentType
	receiverParameters   config.EmailReceiverParameters
	connectionParameters config.EmailConnectionParameters
}

func NewEmail(
	client *mail.Client,
	contentType mail.ContentType,
	receiverParameters config.EmailReceiverParameters,
	connectionParameters config.EmailConnectionParameters,
) *Email {
	return &Email{
		client:               client,
		contentType:          contentType,
		receiverParameters:   receiverParameters,
		connectionParameters: connectionParameters,
	}
}

func (e *Email) Notify(ctx context.Context, rc aggregator.RateChange) error {
	log.Println("Sending email")

	msg, err := e.getMessage(rc.String())

	if err != nil {
		return err
	}

	return e.client.DialAndSendWithContext(ctx, msg)
}

func (e *Email) getMessage(body string) (*mail.Msg, error) {
	msg := mail.NewMsg()

	if err := msg.From(e.connectionParameters.Username); err != nil {
		return nil, err
	}

	if err := msg.To(e.receiverParameters.Email); err != nil {
		return nil, err
	}
	t, err := template.New("tmpl").Parse(bodyTemplate)

	if err != nil {
		return nil, err
	}

	if err := msg.SetBodyHTMLTemplate(t, body); err != nil {
		return nil, err
	}

	msg.Subject(aggregator.NotificationTitle)
	msg.SetBodyString(e.contentType, body)

	return msg, nil
}
