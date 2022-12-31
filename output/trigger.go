package output

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dayvillefire/pocsag-monitor/obj"
)

func init() {
	RegisterOutput("trigger", func() Output { return &TriggerOutput{} })
}

// TriggerOutput pushes a message and cap code to a trigger URL. It makes a
// number of substitutions, including:
//
// * `{cap}` - CAP code
// * `{msg}` - Raw message
//
type TriggerOutput struct {
	url string
}

func (s *TriggerOutput) Init(v string) error {
	s.url = v
	return nil
}

func (s *TriggerOutput) SendMessage(a obj.AlphaMessage, channel, msg string) (string, error) {
	c := s.httpClient()
	u := s.url
	u = strings.ReplaceAll(s.url, "{msg}", url.PathEscape(a.Message))
	u = strings.ReplaceAll(s.url, "{cap}", url.PathEscape(a.CapCode))
	r, err := c.Get(u)
	if err != nil {
		log.Printf("ERR: %s: %s", u, err.Error())
		return "", err
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERR: %s: %s", u, err.Error())
		return "", err
	}
	return string(b), nil
}

func (s *TriggerOutput) httpClient() *http.Client {
	c := http.DefaultClient
	c.Timeout = time.Second * 5
	return c
}
