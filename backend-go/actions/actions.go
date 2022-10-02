package actions

import (
	"time"
)

type Action interface {
	Run() error
	UpdateLastPerformed()
}

type ActionData struct {
	ID            uint
	Enabled       bool
	LastPerformed time.Time
	Message       string
	Data          string
}

func (a *ActionData) UpdateLastPerformed() {
	//entry, err := client.Collection("Action").Doc(ip).Get(c.Request.Context())
}

func (a *ActionData) Run() error {
	return nil
}
