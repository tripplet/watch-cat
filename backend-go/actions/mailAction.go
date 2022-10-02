package actions

import (
	"context"

	"google.golang.org/appengine/mail"
)

type MailAction struct {
	ActionData
	Address string
	Subject string
}

func (a *MailAction) Run() {
	if !a.Enabled {
		return
	}

	a.UpdateLastPerformed()

	msg := &mail.Message{
		Sender:  "event@", // + projectID + ".appspotmail.com", TODO
		To:      []string{a.Address},
		Subject: a.Subject,
		Body:    a.Message,
	}

	mail.Send(context.Background(), msg)
}
