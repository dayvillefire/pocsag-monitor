package output

import (
	"testing"
	"time"

	"github.com/dayvillefire/pocsag-monitor/obj"
)

const (
	uri     = ""
	channel = ""
)

func Test_MQTT(t *testing.T) {
	d := MQTTOutput{}

	err := d.Init(uri)
	if err != nil {
		t.Fatalf(err.Error())
	}
	_, err = d.SendMessage(obj.AlphaMessage{}, channel, "Test message "+time.Now().Format(time.RFC1123))
	if err != nil {
		t.Fatalf(err.Error())
	}
}
