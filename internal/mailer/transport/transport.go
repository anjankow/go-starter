package transport

import "github.com/jordan-wright/email"

//go:generate mockgen -destination=gomock.go -package=transport -mock_names=MailTransporter=GomockMailTransport allaboutapps.dev/aw/go-starter/internal/mailer/transport MailTransporter
type MailTransporter interface {
	Send(mail *email.Email) error
}
