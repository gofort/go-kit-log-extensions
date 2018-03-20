package extensions

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
)

func Test_MultiLoggerEnchanced_Normal(t *testing.T) {

	var buffer bytes.Buffer

	writer := log.NewSyncWriter(&buffer)

	logger := NewMultiLoggerEnchanced(true, time.Minute, log.NewLogfmtLogger(writer), log.NewLogfmtLogger(writer), log.NewLogfmtLogger(os.Stdout))

	err := logger.Log("test", "log")
	if err != nil {
		t.Error(err)
	}

	if buffer.String() != "test=log\n"+"test=log\n" {
		t.Error("unexpected result")
	}

}

func Test_MultiLoggerEnchanced_Timeout(t *testing.T) {

	writer := log.NewSyncWriter(os.Stderr)

	logger := NewMultiLoggerEnchanced(true, time.Nanosecond, log.NewLogfmtLogger(writer),
		log.NewLogfmtLogger(writer), log.NewLogfmtLogger(os.Stdout),
		log.LoggerFunc(func(keyvals ...interface{}) error {
			time.Sleep(time.Second)
			return nil
		}))

	err := logger.Log("test", "log")
	if err == nil {
		t.Error("no timeout!")
	}

}
