package extensions

import (
	"log"

	kitlog "github.com/go-kit/kit/log"
)

// MultiLoggerAsync is a multi logger implementation which call every logger async.
type MultiLoggerAsync struct {
	loggers     []kitlog.Logger
	printErrors bool
}

// NewMultiLoggerAsync creates new MultiLoggerAsync instance.
// If printErrors is true, then MultiLoggerAsync will print all errors of it's child loggers to stderr.
// Warning: if printErrors is false, you will not ever see any errors and they won't be returned in any methods.
func NewMultiLoggerAsync(printErrors bool, loggers ...kitlog.Logger) kitlog.Logger {

	return &MultiLoggerAsync{
		printErrors: printErrors,

		loggers: loggers,
	}

}

// Log implements go-kit logging interface.
func (self *MultiLoggerAsync) Log(keyvals ...interface{}) error {

	for _, logger := range self.loggers {
		logger := logger

		go func(logger kitlog.Logger) {

			err := logger.Log(keyvals...)
			if err != nil && self.printErrors {
				log.Println(err)
			}

		}(logger)
	}

	return nil
}
