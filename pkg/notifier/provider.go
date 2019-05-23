package notifier

import (
	"fmt"
	"github.com/go-mail/mail"
	"github.com/krzysztof-gzocha/curnot/pkg/config"
)

type NotifierProvider struct {
	notifierConfig map[string]config.NotifierConfig
}

func NewNotifierProvider(notifierConfig map[string]config.NotifierConfig) *NotifierProvider {
	return &NotifierProvider{notifierConfig: notifierConfig}
}

func (p *NotifierProvider) GetNotifiers() []NotifierInterface {
	var notifiers []NotifierInterface

	emailNotifierConfig, exists := p.notifierConfig[NameEmailNotifier]

	if exists {
		dialer := mail.NewDialer(
			emailNotifierConfig.ConnectionParameters.Host,
			emailNotifierConfig.ConnectionParameters.Port,
			emailNotifierConfig.ConnectionParameters.Username,
			emailNotifierConfig.ConnectionParameters.Password,
		)

		notifier := NewEmail(
			dialer,
			&emailNotifierConfig.EmailReceiverParameters,
			&emailNotifierConfig.ConnectionParameters,
		)

		notifiers = append(notifiers, notifier)
	}

	desktopNotifier, exists := p.notifierConfig[NameDesktopNotifier]

	if exists {
		fmt.Println(desktopNotifier)

		notifiers = append(notifiers, NewDesktop())
	}

	return notifiers
}
