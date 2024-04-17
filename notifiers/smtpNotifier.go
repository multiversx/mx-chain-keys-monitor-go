package notifiers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"multiversx/mvx-keys-monitor/core"
)

const mimeHeaders = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
const htmlLineBreak = "<br>"
const htmlTemplate = `<!-- template.html -->
<!DOCTYPE html>
<html lang="en">
<body>
   {{.Body}}
</body>
</html>
`

type sendMailHandler func(host string, auth smtp.Auth, from string, to []string, msgBytes []byte) error

type smtpNotifier struct {
	to       string
	smtpPort int
	smtpHost string
	from     string
	password string
	sendMail sendMailHandler
}

// ArgsSmtpNotifier represents the SMTP notifier arguments used in the constructor function
type ArgsSmtpNotifier struct {
	To       string
	SmtpPort int
	SmtpHost string
	From     string
	Password string
}

// NewSmtpNotifier creates a new SMTP email notifier
func NewSmtpNotifier(args ArgsSmtpNotifier) *smtpNotifier {
	return &smtpNotifier{
		to:       args.To,
		smtpPort: args.SmtpPort,
		smtpHost: args.SmtpHost,
		from:     args.From,
		password: args.Password,
		sendMail: sendMail,
	}
}

// OutputMessages will push the provided messages as error
func (notifier *smtpNotifier) OutputMessages(messages ...core.OutputMessage) {
	log.Debug("notifier.OutputMessage pushing error notification as SMTP email", "num messages", len(messages))
	if len(messages) == 0 {
		return
	}

	msgString := ""
	maxMessageOutputType := core.MessageOutputType(0)
	for _, msg := range messages {
		if msg.Type > maxMessageOutputType {
			maxMessageOutputType = msg.Type
		}

		msgString += createMessageString(msg) + htmlLineBreak
	}

	title := createTitle(maxMessageOutputType, messages[0].ExecutorName)

	err := notifier.pushNotification(msgString, title)
	if err != nil {
		log.Error("notifier.OutputMessage pushing notification", "error", err)
	}
}

func (notifier *smtpNotifier) pushNotification(msgString string, title string) error {
	auth := smtp.PlainAuth("", notifier.from, notifier.password, notifier.smtpHost)

	msgBytes, err := createEmailBytes(msgString, title)
	if err != nil {
		return err
	}

	err = notifier.sendMail(
		fmt.Sprintf("%s:%d", notifier.smtpHost, notifier.smtpPort),
		auth,
		notifier.from,
		[]string{notifier.to},
		msgBytes,
	)
	if err != nil {
		return err
	}

	log.Debug("notifier.pushNotification: sent notification as smtp email")

	return nil
}

func sendMail(host string, auth smtp.Auth, from string, to []string, msgBytes []byte) error {
	return smtp.SendMail(host, auth, from, to, msgBytes)
}

func createEmailBytes(msg string, title string) ([]byte, error) {
	var body bytes.Buffer

	mailTemplate := template.New("")
	_, err := mailTemplate.Parse(htmlTemplate)
	if err != nil {
		return nil, err
	}
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", title, mimeHeaders)))

	err = mailTemplate.Execute(&body, struct {
		Body template.HTML
	}{
		Body: template.HTML(msg),
	})
	if err != nil {
		return nil, err
	}

	return body.Bytes(), nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (notifier *smtpNotifier) IsInterfaceNil() bool {
	return notifier == nil
}
