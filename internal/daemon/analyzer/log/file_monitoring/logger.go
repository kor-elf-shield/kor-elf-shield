package file_monitoring

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type logger struct {
	logger log.Logger
}

func NewLogger(log log.Logger) Logger {
	return &logger{logger: log}
}

func (l *logger) Fatal(v ...interface{}) {
	l.logger.Error(fmt.Sprintf("File Monitoring: %v", v...))
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.logger.Error(fmt.Sprintf("File Monitoring: "+format, v...))
}

func (l *logger) Fatalln(v ...interface{}) {
	l.logger.Error(fmt.Sprintf("File Monitoring: %v", v...))
}

func (l *logger) Panic(v ...interface{}) {
	l.logger.Error(fmt.Sprintf("File Monitoring: %v", v...))
}

func (l *logger) Panicf(format string, v ...interface{}) {
	l.logger.Warn(fmt.Sprintf("File Monitoring: "+format, v...))
}

func (l *logger) Panicln(v ...interface{}) {
	l.logger.Error(fmt.Sprintf("File Monitoring: %v", v...))
}

func (l *logger) Print(v ...interface{}) {
	l.logger.Warn(fmt.Sprintf("File Monitoring: %v", v...))
}

func (l *logger) Printf(format string, v ...interface{}) {
	l.logger.Warn(fmt.Sprintf("File Monitoring: "+format, v...))
}

func (l *logger) Println(v ...interface{}) {
	l.logger.Warn(fmt.Sprintf("File Monitoring: %v", v...))
}
