package filter

import (
	"testing"

	"github.com/dayvillefire/pocsag-monitor/obj"
)

func TestFilter(t *testing.T) {
	rf := &RegexFilter{}
	err := rf.Init()
	if err != nil {
		t.Errorf("Init(): ERR: %s", err.Error())
	}
	err = rf.Configure(map[string]any{
		"regexps": []string{
			"ALPHA | ([^|]+)|",
		},
	})
	if err != nil {
		t.Errorf("Configure(): ERR: %s", err.Error())
	}

	o := obj.AlphaMessage{
		Message: "01-03-2025 08:27:35 | AMBULANCE EMERGENCY | SMITH, JOHN & SMYTHE, JANE | 431 KNOB MANOR RD MONTVILLE, 06382 | MFD, M101 | 37 YOF STOMACH PAIN   |",
	}

	o, err = rf.Filter(o)
	if err != nil {
		t.Errorf("Filter(): ERR: %s", err.Error())
	}

	t.Logf("o = %s", o.Message)
}
