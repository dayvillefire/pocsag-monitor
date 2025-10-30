package filter

import (
	"testing"

	"github.com/dayvillefire/pocsag-monitor/obj"
)

func TestSplitMatchFilter(t *testing.T) {
	rf := &SplitMatchFilter{}
	err := rf.Init()
	if err != nil {
		t.Errorf("Init(): ERR: %s", err.Error())
	}
	err = rf.Configure(map[string]any{
		"if-match":      "AMBULANCE EMERGENCY |",
		"split":         "|",
		"remove-fields": []int{2},
	})
	if err != nil {
		t.Errorf("Configure(): ERR: %s", err.Error())
	}

	tests := []obj.AlphaMessage{
		{
			Message: "01-03-2025 08:27:35 | AMBULANCE EMERGENCY | SMITH, JOHN & SMYTHE, JANE | 431 KNOB MANOR RD MONTVILLE, 06382 | MFD, M101 | 37 YOF STOMACH PAIN   |",
		},
		{
			Message: "01-03-2025 08:27:35 | FIRE ALARM | SMITH, JOHN & SMYTHE, JANE | 431 KNOB MANOR RD MONTVILLE, 06382 | MFD, M101 | 37 YOF STOMACH PAIN   |",
		},
	}

	for _, o := range tests {
		o, err = rf.Filter(o)
		if err != nil {
			t.Errorf("Filter(): ERR: %s", err.Error())
		}

		t.Logf("o = %s", o.Message)
	}
}
