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

	logger, err := NewMultiLoggerEnchanced(true, time.Minute, log.NewLogfmtLogger(writer), log.NewLogfmtLogger(writer), log.NewLogfmtLogger(os.Stdout))
	if err != nil {
		t.Error(err)
		return
	}

	err = logger.Log("test", "log")
	if err != nil {
		t.Error(err)
		return
	}

	if buffer.String() != "test=log\n"+"test=log\n" {
		t.Error("unexpected result")
		return
	}

}

func Test_MultiLoggerEnchanced_Timeout(t *testing.T) {

	writer := log.NewSyncWriter(os.Stderr)

	logger, err := NewMultiLoggerEnchanced(true, time.Nanosecond, log.NewLogfmtLogger(writer),
		log.NewLogfmtLogger(writer), log.NewLogfmtLogger(os.Stdout),
		log.LoggerFunc(func(keyvals ...interface{}) error {
			time.Sleep(time.Second)
			return nil
		}))
	if err != nil {
		t.Error(err)
		return
	}

	err = logger.Log("test", "log")
	if err == nil {
		t.Error("no timeout!")
		return
	}

}
