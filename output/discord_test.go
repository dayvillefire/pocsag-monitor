package output

import (
	"testing"

	"github.com/dayvillefire/pocsag-monitor/obj"
)

var (
	discordToken   = ""
	discordChannel = ""
)

func Test_Discord(t *testing.T) {
	d := DiscordOutput{}

	err := d.Init(discordToken)
	if err != nil {
		t.Fatalf(err.Error())
	}
	_, err = d.SendMessage(obj.AlphaMessage{}, discordChannel, "Test message")
	if err != nil {
		t.Fatalf(err.Error())
	}
}
