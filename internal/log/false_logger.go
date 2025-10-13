package log

import "os"

type falseLogger struct{}

func (l *falseLogger) Debug(msg string) {}
func (l *falseLogger) Info(msg string)  {}
func (l *falseLogger) Warn(msg string)  {}
func (l *falseLogger) Error(msg string) {}

func (l *falseLogger) Fatal(msg string) {
	os.Exit(1)
}

func (l *falseLogger) Sync() error { return nil }
