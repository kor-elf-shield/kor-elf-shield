package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	Sync() error
}

type logger struct {
	zap *zap.Logger
	dev bool
}

func (l *logger) Debug(msg string) {
	l.zap.Debug(msg)
}

func (l *logger) Info(msg string) {
	l.zap.Info(msg)
}

func (l *logger) Warn(msg string) {
	l.zap.Warn(msg)
}

func (l *logger) Error(msg string) {
	l.zap.Error(msg)
}

func (l *logger) Fatal(msg string) {
	l.zap.Fatal(msg)
}

func (l *logger) Sync() error {
	return l.zap.Sync()
}

func NewLogger(opts *LoggerOptions) (Logger, error) {
	if !opts.Enabled {
		return &falseLogger{}, nil
	}

	cfg := zap.Config{
		Level:       opts.Level,
		Development: opts.Development,

		// Reduce the amount of noise in logs during mass events.
		Sampling: &zap.SamplingConfig{
			Initial:    100, // The first 100 identical messages per second will be written to the log
			Thereafter: 100, // Every 100th message will be written to the log
		},

		Encoding:         opts.Encoding.String(),
		EncoderConfig:    encoderConfig(opts),
		OutputPaths:      opts.Paths,
		ErrorOutputPaths: opts.LogErrorPaths,
	}
	cfg.DisableStacktrace = !opts.Development
	cfg.DisableCaller = !opts.Development

	zapLogger := zap.Must(cfg.Build())
	return &logger{
		zap: zapLogger,
		dev: opts.Development,
	}, nil
}

func encoderConfig(opts *LoggerOptions) zapcore.EncoderConfig {
	var encoderConfig zapcore.EncoderConfig

	encoderConfig = zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.SkipLineEnding = false
	encoderConfig.TimeKey = "date"

	if opts.Development {
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	}

	return encoderConfig
}
