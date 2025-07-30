package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/dayvillefire/groupme"
	"github.com/dayvillefire/pocsag-monitor/obj"
)

func init() {
	RegisterOutput("groupme", func() Output { return &GroupMeOutput{} })
}

type GroupMeOutput struct {
	botID  string
	token  string
	client groupme.Client
	bot    groupme.Bot
}

func (s *GroupMeOutput) Init(token string) error {
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return fmt.Errorf("bad token -- needs to be botID:token")
	}
	s.token = parts[1]
	s.botID = parts[0]
	return nil
}

func (s *GroupMeOutput) SendMessage(a obj.AlphaMessage, channel, msg string) (string, error) {
	s.client = groupme.NewClient(groupme.V3BaseURL, s.token)
	s.bot = groupme.NewBot(groupme.V3BaseURL, s.botID, channel, "", "")

	_, err := s.client.CreateMessage(channel, fmt.Sprintf("%s", time.Now().String()), msg)
	//log.Printf("ret = %#v", ret)
	return "", err
}
