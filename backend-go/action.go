package main

import "time"

type action interface {
	Run() error
	UpdateLastPerformed()
}

type actionDoc struct {
	Enabled         bool      `firestore:"enabled"`
	LastPerformed   time.Time `firestore:"last_performed"`
	IsFailureAction bool      `firestore:"failure_action"`
	Message         string    `firestore:"message"`
}

func (a *actionDoc) UpdateLastPerformed() {
	//entry, err := client.Collection("Action").Doc(ip).Get(c.Request.Context())
}
