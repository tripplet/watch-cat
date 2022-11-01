package models

import (
	"time"
	"watchcat/env"
)

type LogEntry struct {
	Job       uint `gorm:"foreignkey:JobID"`
	Name      string
	EventTime time.Time
	Detail    string
}

func CreateLogEntry(env *env.Env, job uint, name string, detail string) {
	entry := LogEntry{
		Job:       job,
		Name:      name,
		EventTime: time.Now(),
		Detail:    detail,
	}

	env.Database.Create(&entry)
}
