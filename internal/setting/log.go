package setting

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	log2 "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/filesystem"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
)

type log struct {
	Enabled       bool     `mapstructure:"enabled"`
	Level         string   `mapstructure:"level"`
	Development   bool     `mapstructure:"development"`
	Encoding      string   `mapstructure:"encoding"`
	Paths         []string `mapstructure:"paths"`
	LogErrorPaths []string `mapstructure:"log_error_paths"`
}

func logDefault() *log {
	return &log{
		Enabled:       true,
		Level:         "info",
		Development:   false,
		Encoding:      "json",
		Paths:         []string{"/var/log/kor-elf-shield.log"},
		LogErrorPaths: []string{"stderr"},
	}
}

func (l *log) ToLoggerOptions() (log2.LoggerOptions, error) {
	var err error

	level, err := log2.ParseLevel(l.Level)
	if err != nil {
		return log2.LoggerOptions{}, err
	}

	encoding, err := parseEncoding(l.Encoding)
	if err != nil {
		return log2.LoggerOptions{}, err
	}

	paths, err := l.paths()
	if err != nil {
		return log2.LoggerOptions{}, err
	}

	logErrorPaths, err := l.logErrorPaths()
	if err != nil {
		return log2.LoggerOptions{}, err
	}

	return log2.LoggerOptions{
		Enabled:       l.Enabled,
		Level:         level,
		Development:   l.Development,
		Encoding:      encoding,
		Paths:         paths,
		LogErrorPaths: logErrorPaths,
	}, nil
}

func (l *log) paths() ([]string, error) {
	if !l.Enabled {
		return []string{}, nil
	}

	if len(l.Paths) == 0 {
		return []string{}, errors.New(i18n.Lang.T("parameter is not specified", map[string]any{
			"Parameter": "paths",
		}))
	}

	if err := validatePaths(l.Paths); err != nil {
		return []string{}, err
	}

	return l.Paths, nil
}

func (l *log) logErrorPaths() ([]string, error) {
	if !l.Enabled {
		return []string{}, nil
	}

	if len(l.LogErrorPaths) == 0 {
		return []string{}, errors.New(i18n.Lang.T("parameter is not specified", map[string]any{
			"Parameter": "log_error_paths",
		}))
	}

	if err := validatePaths(l.LogErrorPaths); err != nil {
		return []string{}, err
	}

	return l.LogErrorPaths, nil
}

func parseEncoding(encoding string) (log2.Encoding, error) {
	switch encoding {
	case "json":
		return log2.JsonEncoding, nil
	case "text":
		return log2.ConsoleEncoding, nil
	default:
		return log2.JsonEncoding, errors.New(i18n.Lang.T("invalid log encoding", map[string]any{
			"Encoding":  encoding,
			"Encodings": "json, text",
		}))
	}
}

func validatePaths(paths []string) error {
	for _, path := range paths {
		if path == "stdout" || path == "stderr" {
			continue
		}

		err := validatePath(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func validatePath(path string) error {
	if err := validate.PathFile(path, "log.paths"); err != nil {
		return err
	}

	if !strings.HasSuffix(strings.ToLower(path), ".log") {
		return fmt.Errorf("invalid %s. Must be .log", "log.paths")
	}

	dir := filepath.Dir(path)
	if err := filesystem.EnsureDir(dir); err != nil {
		return err
	}

	if err := filesystem.FileHasWritePermissions(path); err != nil {
		return err
	}

	return nil
}
