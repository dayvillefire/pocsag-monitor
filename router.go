package main

import (
	"strings"

	"github.com/dayvillefire/pocsag-monitor/obj"
	"github.com/jbuchbinder/shims"
)

type Router struct {
	ChannelMappings map[string][]string
}

func (r *Router) MapMessage(msg obj.AlphaMessage) []string {
	out := []string{}

	if len(msg.Message) < 4 {
		// Don't post messages shorter than "TEST"
		return out
	}

	found := false

	for k, v := range r.ChannelMappings {
		if strings.HasPrefix(k, "S:") {
			if matchMessage(k[2:], msg.Message) {
				found = true
				out = append(out, v...)
			}
			continue
		}

		if matchCap(k, msg.CapCode) {
			found = true
			out = append(out, v...)
		}
	}

	// Check for a "DEFAULT" mapping
	if !found {
		if m, ok := r.ChannelMappings["DEFAULT"]; ok {
			out = append(out, m...)
		}
	}

	return shims.UniqueValues(out)
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

func matchMessage(mapping, msg string) bool {
	return strings.Contains(msg, mapping)
}
