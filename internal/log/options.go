package log

import (
	"errors"
	"fmt"
	"kor-elf-shield/internal/i18n"

	"go.uber.org/zap"
)

type Encoding int8

const (
	JsonEncoding Encoding = iota + 1
	ConsoleEncoding
)

type LoggerOptions struct {
	Enabled       bool
	Level         zap.AtomicLevel
	Development   bool
	Encoding      Encoding
	Paths         []string
	LogErrorPaths []string
}

func ParseLevel(level string) (zap.AtomicLevel, error) {
	atomicLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		errMessage := i18n.Lang.T(
			"invalid log level",
			map[string]any{
				"Level":  level,
				"Levels": "debug, info, warn, error, fatal",
			},
		)
		return atomicLevel, errors.New(errMessage)
	}

	return atomicLevel, nil
}

func (e Encoding) String() string {
	switch e {
	case JsonEncoding:
		return "json"
	case ConsoleEncoding:
		return "console"
	default:
		return fmt.Sprintf("Encoding(%d)", e)
	}
}
