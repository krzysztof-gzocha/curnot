package notifier

import (
	"context"
	"fmt"

	"github.com/go-mail/mail"
	"github.com/krzysztof-gzocha/curnot/pkg/aggregator"
	"github.com/krzysztof-gzocha/curnot/pkg/config"
	"github.com/krzysztof-gzocha/curnot/pkg/formatter"
)

type Email struct {
	dialer               mail.Dialer
	receiverParameters   config.EmailReceiverParameters
	connectionParameters config.EmailConnectionParameters
	formatter            formatter.MessageFormatStrategy
}

func NewEmail(
	dialer mail.Dialer,
	receiverParameters config.EmailReceiverParameters,
	connectionParameters config.EmailConnectionParameters,
	formatter formatter.MessageFormatStrategy) *Email {
	return &Email{
		dialer:               dialer,
		receiverParameters:   receiverParameters,
		connectionParameters: connectionParameters,
		formatter:            formatter,
	}
}

func (e *Email) Notify(_ context.Context, msg aggregator.RateChange) error {
	fmt.Println("Sending email")

	return e.dialer.DialAndSend(e.getMessage(msg.String()))
}

func (e *Email) getMessage(body string) *mail.Message {
	message := mail.NewMessage()

	message.SetHeader("From", e.connectionParameters.Username)
	message.SetHeader("To", e.receiverParameters.Email)
	message.SetHeader("Subject", aggregator.NotificationTitle)
	message.SetBody(e.formatter.Format(body))

	return message
}
