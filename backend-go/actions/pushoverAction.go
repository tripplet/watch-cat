package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PushoverAction struct {
	ActionData
	IsFailureAction bool
	Token           string
	UserKeys        []string
	CustomSound     string
}

type pushoverApiData struct {
	Token   string `json:"token"`
	User    string `json:"user"`
	Message string `json:"message"`
	Sound   string `json:"sound"`
}

func (a *PushoverAction) Run() error {
	if !a.Enabled {
		return nil
	}

	a.UpdateLastPerformed()

	var sound string
	if a.CustomSound != "" {
		sound = a.CustomSound
	} else {
		if a.IsFailureAction {
			sound = "falling"
		} else {
			sound = "bike"
		}
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var errors []error
	for _, user := range a.UserKeys {
		data, err := json.Marshal(pushoverApiData{
			Token:   a.Token,
			User:    user,
			Message: a.Message,
			Sound:   sound})

		if err != nil {
			errors = append(errors, err)
			continue
		}

		// TODO use goroutine
		response, err := client.Post("https://api.pushover.net/1/messages.json", "application/json", bytes.NewBuffer(data))

		if err != nil {
			errors = append(errors, err)
			continue
		}

		if response.StatusCode != http.StatusOK {
			errors = append(errors, fmt.Errorf("pushover service responded with %d", response.StatusCode))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("pushover service responded with %d errors", len(errors))
	}

	return nil
}
