package main

import "time"

type LogEntry struct {
	Job       int64
	Name      string
	EventTime time.Time
	Detail    string
}
