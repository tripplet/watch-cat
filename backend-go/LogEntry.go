package main

import "time"

type LogEntry struct {
	Job       uint `gorm:"foreignkey:JobID"`
	Name      string
	EventTime time.Time
	Detail    string
}

func (env *Env) createLogEntry(job uint, name string, detail string) {
	entry := LogEntry{
		Job:       job,
		Name:      name,
		EventTime: time.Now(),
		Detail:    detail,
	}

	env.db.Create(&entry)
}
