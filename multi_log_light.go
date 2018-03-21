package extensions

import (
	"github.com/go-kit/kit/log"
)

type MultiLogLight struct {
	loggers []log.Logger
}

func (self *MultiLogLight) Log(keyvals ...interface{}) error {
	for _, logger := range self.loggers {
		logger.Log(keyvals...)
	}
	return nil
}

func MultiLogWrapper(loggers ...log.Logger) log.Logger {
	return &MultiLogLight{loggers: loggers}
}
