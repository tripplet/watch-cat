package main

import (
	"time"
)

type action interface {
	Run() error
	UpdateLastPerformed()
}

type actionData struct {
	ID            uint
	Enabled       bool
	LastPerformed time.Time
	Message       string
	Data          string
}

func (a *actionData) UpdateLastPerformed() {
	//entry, err := client.Collection("Action").Doc(ip).Get(c.Request.Context())
}
