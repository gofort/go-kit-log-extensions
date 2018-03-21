package extensions

import (
	"testing"
)

func Test_SlackBackend(t *testing.T) {

	slackLogger, err := NewSlackBackend("", 0, "localhost", "test_component")
	if err != nil {
		t.Error(err)
		return
	}

	slackLogger.Log("structured", "loggging")

}
