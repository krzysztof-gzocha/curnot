package notifier

import (
	"fmt"
	"github.com/go-mail/mail"
	"github.com/krzysztof-gzocha/curnot/pkg/config"
)

const NameEmailNotifier = "email"

type Email struct {
	dialer               *mail.Dialer
	receiverParameters   *config.EmailReceiverParameters
	connectionParameters *config.EmailConnectionParameters
}

func NewEmail(
	dialer *mail.Dialer,
	receiverParameters *config.EmailReceiverParameters,
	connectionParameters *config.EmailConnectionParameters) *Email {
	return &Email{
		dialer:               dialer,
		receiverParameters:   receiverParameters,
		connectionParameters: connectionParameters,
	}
}

func (e *Email) Notify(msg string) error {
	fmt.Println("Sending email")

	if e.dialer == nil {
		panic(fmt.Sprintf("Dialer hasn't been set"))
	}

	message := e.getMessage(msg)

	if err := e.dialer.DialAndSend(message); err != nil {
		panic(err)
	}

	return nil
}

func (e *Email) getMessage(body string) *mail.Message {
	message := mail.NewMessage()

	message.SetHeader("From", e.connectionParameters.Username)
	message.SetHeader("To", e.receiverParameters.Email)
	message.SetHeader("Subject", NotificationTitle)
	message.SetBody("text/html", fmt.Sprintf(body))

	return message
}
