package log

import "time"

type Entry struct {
	Message string
	Unit    string
	PID     string
	Time    time.Time
}
