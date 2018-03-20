package extensions

import (
	"log"
	"time"

	kitlog "github.com/go-kit/kit/log"
	multierror "github.com/hashicorp/go-multierror"
)

// MultiLogger is a multi logger implementation which call every logger consistently and has timeout per Log call (not for all loggers separately).
type MultiLogger struct {
	printErrors bool
	timeout     time.Duration

	loggers []kitlog.Logger
}

// NewMultiLogger creates new MultiLogger instance.
// If printErrors is true, then besides returning errors within Log function
// MultiLogger will print all errors of it's child loggers to stderr.
// To exclude situations when one of the loggers will hang up, timeout variable exists.
func NewMultiLogger(printErrors bool, timeout time.Duration, loggers ...kitlog.Logger) kitlog.Logger {

	return &MultiLogger{
		printErrors: printErrors,
		timeout:     timeout,

		loggers: loggers,
	}

}

// Log implements go-kit logging interface.
func (self *MultiLogger) Log(keyvals ...interface{}) error {

	if self.timeout < 1 {
		return self.log(keyvals...)
	}

	finished := make(chan error)

	go func() {
		finished <- self.log(keyvals...)
	}()

	select {
	case err := <-finished:
		close(finished)
		return err
	case <-time.After(self.timeout):
		return ErrTimeout
	}

}

func (self *MultiLogger) log(keyvals ...interface{}) error {

	multiErr := new(multierror.Error)

	for _, logger := range self.loggers {
		err := logger.Log(keyvals...)

		if err != nil {
			multiErr = multierror.Append(multiErr, err)

			if self.printErrors {
				log.Println(err)
			}
		}

	}

	return multiErr.ErrorOrNil()
}
