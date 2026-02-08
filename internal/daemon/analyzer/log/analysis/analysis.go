package analysis

import (
	"errors"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
)

type Entry struct {
	Source           config.SourceType
	Message          string
	Unit             string
	PID              string
	SyslogIdentifier string
	Time             time.Time
}

type regexField struct {
	name  string
	value string
}

func getValueStartEndByRegexIndex(valueId int, idx []int) (start int, end int, err error) {
	id := 2 * valueId

	if idx == nil || len(idx) <= id+1 {
		return 0, 0, errors.New("invalid index")
	}

	start, end = idx[id], idx[id+1]
	if start < 0 || end < 0 {
		return 0, 0, errors.New("invalid index")
	}

	return start, end, nil
}
