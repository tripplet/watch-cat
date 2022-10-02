package main

import (
	"context"

	"google.golang.org/appengine/mail"
)

type mailAction struct {
	actionData
	Address string
	Subject string
}

func (a *mailAction) Run() {
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
