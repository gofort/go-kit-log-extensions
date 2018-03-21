package extensions

import (
	"github.com/go-kit/kit/log"
)

// DoubleLog offers you ability to log using two loggers
func DoubleLog(stdlog log.Logger, additionalLog log.Logger, keyvals ...interface{}) {
	stdlog.Log(keyvals...)
	additionalLog.Log(keyvals...)
}
