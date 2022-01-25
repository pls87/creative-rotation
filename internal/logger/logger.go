package logger

import (
	"os"

	"github.com/pls87/creative-rotation/internal/config"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func New(cfg config.LoggerConf) *Logger {
	log := &Logger{logrus.New()}
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel
	if lvl, err := logrus.ParseLevel(cfg.Level); err != nil {
		log.Level = lvl
	}
	return log
}

func (l *Logger) IfError(e error, msg string) {
	if e != nil {
		l.Errorf("%s: %s", msg, e)
	}
}
