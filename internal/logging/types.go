package logging

import "context"

type Configuration interface {
	Level() string
	Name() string
	Format() string
	Caller() bool
	Env() string
}

// Logger is app provided logging interface
type Logger interface {
	WithContext(ctx context.Context) Logger
	WithError(err error) Logger
	WithField(tag Tag, value any) Logger
	Info(msg string)
	Warn(msg string)
	Debug(msg string)
	Error(msg string)
	Fatal(msg string)
	Panic(msg string)
}
