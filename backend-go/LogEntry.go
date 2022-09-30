package main

import "time"

type LogEntry struct {
	Job       uint `gorm:"foreignkey:JobID"`
	Name      string
	EventTime time.Time
	Detail    string
}
