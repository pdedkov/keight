//nolint:ireturn // it's ok here
package logging

import (
	"context"
	"log/slog"
	"maps"
	"os"
)

type loggingKey int

const (
	_ loggingKey = iota
	loggerKey
)

type SLog struct {
	log  *slog.Logger
	ctx  context.Context //nolint:containedctx // it's ok here
	err  error
	data map[Tag]any
}

// NewSLog creates new slog logger
func NewSLog(appName string, cfg Configuration) *SLog {
	opts := &slog.HandlerOptions{
		AddSource: cfg.Caller(),
	}
	switch cfg.Level() {
	case "debug":
		opts.Level = slog.LevelDebug
	case "error", "fatal":
		opts.Level = slog.LevelError
	case "warning":
		opts.Level = slog.LevelWarn
	case "info":
		opts.Level = slog.LevelInfo
	default:
		opts.Level = slog.LevelInfo
	}

	var handler slog.Handler
	switch cfg.Format() {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	if cfg.Name() != "" {
		appName = cfg.Name()
	}
	logger := slog.New(handler.WithAttrs([]slog.Attr{
		slog.String(Env.String(), cfg.Env()), slog.String(App.String(), appName),
	}))
	slog.SetDefault(logger)

	return &SLog{
		log:  logger,
		data: make(map[Tag]any),
	}
}

// WithContext adds context to record
func (s *SLog) WithContext(ctx context.Context) Logger {
	c := s.clone()
	c.ctx = ctx
	return c
}

// Info information log message
func (s *SLog) Info(msg string) {
	s.log.InfoContext(s.ctx, msg, s.attrs()...)
}

// Warn warning log message
func (s *SLog) Warn(msg string) {
	s.log.WarnContext(s.ctx, msg, s.attrs()...)
}

// Debug log debug message
func (s *SLog) Debug(msg string) {
	s.log.DebugContext(s.ctx, msg, s.attrs()...)
}

// Error err log message
func (s *SLog) Error(msg string) {
	s.log.ErrorContext(s.ctx, msg, s.attrs()...)
}

// Fatal error message and exit
func (s *SLog) Fatal(msg string) {
	s.log.ErrorContext(s.ctx, msg, s.attrs()...)
	os.Exit(1)
}

// Panic error message and panic
func (s *SLog) Panic(msg string) {
	s.log.ErrorContext(s.ctx, msg, s.attrs()...)
	panic(msg)
}

func (s *SLog) clone() *SLog {
	if s.data == nil {
		s.data = make(map[Tag]any)
	}
	data := maps.Clone(s.data)
	return &SLog{
		err:  s.err,
		ctx:  s.ctx,
		data: data,
		log:  s.log,
	}
}

// WithField adds fields to logging
func (s *SLog) WithField(tag Tag, value any) Logger {
	c := s.clone()
	c.data[tag] = value
	return c
}

func (s *SLog) attrs() []any {
	attributes := make([]any, 0, len(s.data)*2)
	for key, value := range s.data {
		attributes = append(attributes, key.String(), value)
	}
	if s.err != nil {
		attributes = append(attributes, Error.String(), s.err)
	}
	return attributes
}

// WithError adds error to slog
func (s *SLog) WithError(err error) Logger {
	c := s.clone()
	c.err = err
	return c
}

// NewContext add logging service to context
func NewContext(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

// Default creates logger from default slog
func Default() Logger {
	return &SLog{
		log: slog.Default(),
	}
}

// FromContext returns log from context
func FromContext(ctx context.Context) Logger {
	log, _ := ctx.Value(loggerKey).(Logger)
	return log
}
