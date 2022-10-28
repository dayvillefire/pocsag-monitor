package main

import "github.com/dayvillefire/pocsag-monitor/obj"

type Router struct {
	ChannelMappings map[string]string
}

func (r *Router) MapMessage(msg obj.AlphaMessage) []string {
	out := []string{}

	// Check for a "DEFAULT" mapping
	if m, ok := r.ChannelMappings["DEFAULT"]; ok {
		out = append(out, m)
	}

	for k, v := range r.ChannelMappings {
		if matchCap(k, msg.CapCode) {
			out = append(out, v)
		}
	}

	return out
}

func matchCap(mapping, cap string) bool {
	// Make sure it's 7 digits long prefixed by 0's
	mycap := ""
	if len(cap) < 7 {
		for i := 0; i < 7-len(cap); i++ {
			mycap = mycap + "0"
		}
		mycap = mycap + cap
	} else {
		mycap = cap
	}

	// Direct mapping to save a few cycles
	if mapping == mycap {
		return true
	}

	// Byte by byte mapping
	mappingRunes := []rune(mapping)
	capRunes := []rune(mycap)
	for pos, _ := range mappingRunes {
		if mappingRunes[pos] == 'X' || mappingRunes[pos] == 'x' {
			continue
		}
		if capRunes[pos] != mappingRunes[pos] {
			return false
		}
	}
	return true
}
