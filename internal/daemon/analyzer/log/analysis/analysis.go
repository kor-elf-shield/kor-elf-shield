package analysis

import (
	"time"
)

type Analysis interface {
	Process(entry *Entry) error
}

type Entry struct {
	Message string
	Unit    string
	PID     string
	Time    time.Time
}

type EmptyAnalysis struct{}

func (empty *EmptyAnalysis) Process(_ *Entry) error {
	return nil
}
