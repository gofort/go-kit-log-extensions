package extensions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-logfmt/logfmt"
)

// SlackBackend implements log interface with ability to log into slack
type SlackBackend struct {
	client *http.Client
	url    *url.URL

	username string
	author   string
}

// NewSlackBackend creates new slack backend, where urlStr is slack hook url,
// requestTimeout is a timeout per http request, username is message username and
// author is an author of attachement, see https://api.slack.com/docs/message-formatting for more info
func NewSlackBackend(urlStr string, requestTimeout time.Duration, username, author string) (log.Logger, error) {

	if urlStr == "" {
		return nil, errors.New("slack url can't be empty")
	}

	if requestTimeout < 1 {
		requestTimeout = time.Second * 2
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return &SlackBackend{
		url:      u,
		author:   author,
		username: username,
		client: &http.Client{
			Timeout: requestTimeout,
		},
	}, nil

}

// Log implements go-kit logging interface.
func (self *SlackBackend) Log(keyvals ...interface{}) error {

	var buffer bytes.Buffer

	encoder := logfmt.NewEncoder(&buffer)
	err := encoder.EncodeKeyvals(keyvals...)
	if err != nil {
		return err
	}

	err = encoder.EndRecord()
	if err != nil {
		return err
	}

	message := slackPayload{
		Username:     self.username,
		Attachements: []slackAttachement{{Color: "#FF0000", Text: buffer.String(), AuthorName: self.author}},
	}

	content, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", self.url.String(), bytes.NewBuffer(content))
	if err != nil {
		return err
	}
	req.Header.Set("Content-type", "application/json")

	res, err := self.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBodyBytes, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("unable to send message to slack: status_code=%d error=%s", res.StatusCode, string(resBodyBytes))
	}

	return nil
}

type slackPayload struct {
	Username     string             `json:"username"`
	Attachements []slackAttachement `json:"attachments"`
}

type slackAttachement struct {
	Color      string `json:"color"`
	Text       string `json:"text"`
	AuthorName string `json:"author_name"`
}
