package log

import "os"

type falseLogger struct{}

func (l *falseLogger) Debug(_ string) {}
func (l *falseLogger) Info(_ string)  {}
func (l *falseLogger) Warn(_ string)  {}
func (l *falseLogger) Error(_ string) {}

func (l *falseLogger) Fatal(_ string) {
	os.Exit(1)
}

func (l *falseLogger) Sync() error { return nil }
func (l *falseLogger) ReOpen() error {
	return nil
}

func NewFalseLogger() Logger {
	return &falseLogger{}
}
