package config

type Log struct {
	LogName   string `envconfig:"LOG_APP_NAME" default:""`
	LogLevel  string `envconfig:"LOG_LEVEL" default:"info"`
	LogFormat string `envconfig:"LOG_FORMAT" default:"text"`
	LogCaller bool   `envconfig:"LOG_CALLER" default:"false"`
	LogEnv    string `envconfig:"LOG_ENV"`
}

func (l *Log) Name() string {
	return l.LogName
}

// Level is default logging level: debug,info,warning,error,falat
func (l *Log) Level() string {
	return l.LogLevel
}

// Format default logging format, default text (json for json logs)
func (l *Log) Format() string {
	return l.LogFormat
}

// Caller add caller string to log records
func (l *Log) Caller() bool {
	return l.LogCaller
}

// Env app env (local, staging, prod, etc)
func (l *Log) Env() string {
	return l.LogEnv
}
