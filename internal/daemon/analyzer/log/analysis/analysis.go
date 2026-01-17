package analysis

import (
	"time"
)

type Analysis interface {
	Process(entry *Entry) error
}

type Entry struct {
	Message          string
	Unit             string
	PID              string
	SyslogIdentifier string
	Time             time.Time
}

type processReturn struct {
	found   bool
	subject string
	body    string
}

type EmptyAnalysis struct{}

func (empty *EmptyAnalysis) Process(_ *Entry) error {
	return nil
}
