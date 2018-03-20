package extensions

import (
	"log"
	"time"

	kitlog "github.com/go-kit/kit/log"
	multierror "github.com/hashicorp/go-multierror"
)

// MultiLoggerEnchanced is a multi logger implementation which call every logger async,
// but waits until all of them will be finished (configurable timeout per logger exists).
type MultiLoggerEnchanced struct {
	printErrors      bool
	timeoutPerLogger time.Duration

	loggers []kitlog.Logger
}

// NewMultiLoggerEnchanced creates new MultiLoggerEnchanced instance.
// If printErrors is true, then besides returning errors within Log function
// MultiLoggerEnchanced will print all errors of it's child loggers to stderr.
// To exclude situations when one of the loggers will hang up, timeoutPerLogger variable exists.
func NewMultiLoggerEnchanced(printErrors bool, timeoutPerLogger time.Duration, loggers ...kitlog.Logger) kitlog.Logger {

	if timeoutPerLogger < 1 {
		panic("timeout per logger can't be less than 1")
	}

	return &MultiLoggerEnchanced{
		printErrors:      printErrors,
		timeoutPerLogger: timeoutPerLogger,

		loggers: loggers,
	}

}

// Log implements go-kit logging interface.
func (self *MultiLoggerEnchanced) Log(keyvals ...interface{}) error {

	finishedAll := make(chan error, len(self.loggers))

	for _, logger := range self.loggers {
		logger := logger

		go func() {
			finishedAll <- self.logWithTimeout(logger, keyvals...)
		}()

	}

	multiErr := new(multierror.Error)

	returnsNumber := 0

	for {
		err := <-finishedAll

		if err != nil {
			multiErr = multierror.Append(multiErr, err)

			if self.printErrors {
				log.Println(err)
			}
		}

		returnsNumber++

		if returnsNumber == len(self.loggers) {
			break
		}

	}

	close(finishedAll)

	return multiErr.ErrorOrNil()
}

func (self *MultiLoggerEnchanced) logWithTimeout(logger kitlog.Logger, keyvals ...interface{}) error {

	finished := make(chan error)

	go func() {
		finished <- logger.Log(keyvals...)
	}()

	select {
	case err := <-finished:
		close(finished)
		return err
	case <-time.After(self.timeoutPerLogger):
		return ErrTimeout
	}

}
