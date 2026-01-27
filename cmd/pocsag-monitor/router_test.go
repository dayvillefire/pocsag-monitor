package main

import (
	"testing"

	"github.com/dayvillefire/pocsag-monitor/obj"
)

func Test_Router(t *testing.T) {
	r := Router{
		ChannelMappings: map[string][]string{
			"DEFAULT": {"012345"},
			"0620XXX": {"NOT THIS ONE"},
			"0630XXX": {"67890"},
		},
	}
	dest := r.MapMessage(obj.AlphaMessage{CapCode: "0630777"})
	if len(dest) != 2 {
		t.Errorf("Wrong mappings: %#v", dest)
	}
}

func Test_MatchCap(t *testing.T) {
	type testItem struct {
		pattern string
		cap     string
		match   bool
		descrip string
	}
	testSuite := []testItem{
		{"0630XXX", "0630777", true, "basic wildcard match"},
		{"0630XXX", "630777", true, "wildcard match, no prefix"},
	}

	for _, titem := range testSuite {
		if matchCap(titem.pattern, titem.cap) != titem.match {
			if titem.match {
				t.Errorf("%s should match %s (%s)", titem.pattern, titem.cap, titem.descrip)
			} else {
				t.Errorf("%s should not match %s (%s)", titem.pattern, titem.cap, titem.descrip)
			}
		}
	}
}
