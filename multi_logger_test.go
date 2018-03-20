package extensions

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
)

func Test_MultiLogger_Normal(t *testing.T) {

	var buffer bytes.Buffer

	writer := log.NewSyncWriter(&buffer)

	logger := NewMultiLogger(true, time.Minute, log.NewLogfmtLogger(writer), log.NewLogfmtLogger(writer), log.NewLogfmtLogger(os.Stdout))

	err := logger.Log("test", "log")
	if err != nil {
		t.Error(err)
		return
	}

	if buffer.String() != "test=log\n"+"test=log\n" {
		t.Error("unexpected result")
	}

}

func Test_MultiLogger_Timeout(t *testing.T) {

	writer := log.NewSyncWriter(os.Stderr)

	logger := NewMultiLogger(true, time.Nanosecond, log.NewLogfmtLogger(writer),
		log.NewLogfmtLogger(writer), log.NewLogfmtLogger(os.Stdout),
		log.LoggerFunc(func(keyvals ...interface{}) error {
			time.Sleep(time.Second)
			return nil
		}))

	err := logger.Log("test", "log")
	if err == nil {
		t.Error("no timeout!")
		return
	}

}
