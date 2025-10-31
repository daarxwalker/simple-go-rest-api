package logger

import (
	"github.com/goforj/godump"

	"gocourse/pkg/env"
)

type Logger struct {
}

func New() *Logger {
	return &Logger{}
}

func (s *Logger) Error(err error) {
	if env.Production() {
		// TODO: send to sentry. prometheus, etc...
		return
	}
	godump.Dump(err)
}
