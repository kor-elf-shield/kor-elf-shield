package logger_mock

import "sync"

type LoggerMock struct {
	mu sync.Mutex

	DebugMessages []string
	InfoMessages  []string
	WarnMessages  []string
	ErrorMessages []string
	FatalMessages []string

	SyncErr   error
	ReOpenErr error

	SyncCalls   int
	ReOpenCalls int
}

func New() *LoggerMock {
	return &LoggerMock{
		DebugMessages: make([]string, 0),
		InfoMessages:  make([]string, 0),
		WarnMessages:  make([]string, 0),
		ErrorMessages: make([]string, 0),
		FatalMessages: make([]string, 0),
	}
}

func (l *LoggerMock) Debug(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.DebugMessages = append(l.DebugMessages, msg)
}

func (l *LoggerMock) Info(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.InfoMessages = append(l.InfoMessages, msg)
}

func (l *LoggerMock) Warn(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.WarnMessages = append(l.WarnMessages, msg)
}

func (l *LoggerMock) Error(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.ErrorMessages = append(l.ErrorMessages, msg)
}

func (l *LoggerMock) Fatal(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.FatalMessages = append(l.FatalMessages, msg)
}

func (l *LoggerMock) Sync() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.SyncCalls++
	return l.SyncErr
}

func (l *LoggerMock) ReOpen() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.ReOpenCalls++
	return l.ReOpenErr
}
