package geoip

import (
	"git.kor-elf.net/kor-elf-shield/geoip2"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type logger struct {
	logger log.Logger
}

func NewLogger(log log.Logger) geoip2.Logger {
	return &logger{logger: log}
}

func (l *logger) Error(err error) {
	l.logger.Error(err.Error())
}
