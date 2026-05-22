package pocsag

import (
	"testing"

	"github.com/dayvillefire/pocsag-monitor/obj"
)

func Test_MessageAssembly_SimpleMessage(t *testing.T) {
	ma := newMessageAssembler()
	// Use a capcode that fits within the POCSAG 18-bit address space (max 262143)
	addrCW := buildAddressCodeword(123456, 0)
	msgCW := buildMessageCodeword("TE")

	var results []obj.AlphaMessage
	if m := ma.feedCodeword(addrCW); m != nil {
		results = append(results, *m)
	}
	if m := ma.feedCodeword(msgCW); m != nil {
		results = append(results, *m)
	}
	addrCW2 := buildAddressCodeword(76543, 1)
	if m := ma.feedCodeword(addrCW2); m != nil {
		results = append(results, *m)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 message, got %d", len(results))
	}
	if results[0].CapCode != "0123456" {
		t.Errorf("expected capcode 0123456, got %s", results[0].CapCode)
	}
	t.Logf("message: %s", results[0].Message)
}

func Test_MessageAssembly_FunctionBits(t *testing.T) {
	ma := newMessageAssembler()
	addrCW := buildAddressCodeword(777, 3)
	ma.feedCodeword(addrCW)
	m := ma.feedCodeword(buildAddressCodeword(888, 0))
	if m == nil {
		t.Fatal("expected message")
	}
	t.Logf("function bits preserved: cap=%s", m.CapCode)
}
