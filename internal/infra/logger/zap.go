package logger

import (
	"plantao/internal/domain/log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	log    *zap.SugaredLogger
	rawZap *zap.Logger
}

func NewLogger() (log.Logger, error) {
	var cfg zap.Config

	cfg = zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build(zap.AddCallerSkip(1)) // CallerSkip para mostrar o caller correto

	if err != nil {
		return nil, err
	}

	return &zapLogger{
		log:    logger.Sugar(),
		rawZap: logger,
	}, nil
}

func (l *zapLogger) Info(msg string, args ...any) {
	if !l.validateArgs(msg, args) {
		return
	}

	l.log.Infow(msg, args...)
}

func (l *zapLogger) Warn(msg string, args ...any) {
	if !l.validateArgs(msg, args) {
		return
	}

	l.log.Warnw(msg, args...)
}

func (l *zapLogger) Error(msg string, args ...any) {
	if !l.validateArgs(msg, args) {
		return
	}

	l.log.Errorw(msg, args...)
}

func (l *zapLogger) Debug(msg string, args ...any) {
	if !l.validateArgs(msg, args) {
		return
	}

	l.log.Debugw(msg, args...)
}

func (l *zapLogger) Fatal(msg string, args ...any) {
	if !l.validateArgs(msg, args) {
		return
	}

	l.log.Fatalw(msg, args...)
}

func (l *zapLogger) With(args ...any) log.Logger {
	if len(args)%2 != 0 {
		l.log.Warnw("args inválidos passados para o logger, número ímpar de argumentos",
			"args", args,
		)
		return l
	}

	return &zapLogger{
		log:    l.log.With(args...),
		rawZap: l.rawZap,
	}
}

func (l *zapLogger) Sync() error {
	return l.rawZap.Sync()
}

func (l *zapLogger) validateArgs(msg string, args []any) bool {
	if len(args)%2 != 0 {
		l.log.Warnw("args inválidos passados para o logger",
			"msg", msg,
			"args", args,
		)
		return false
	}
	return true
}
