package pocsag

import (
	"strings"
	"time"

	"github.com/dayvillefire/pocsag-monitor/obj"
)

type messageAssembler struct {
	currentCapCode string
	currentFunc    int
	msgBuf         strings.Builder
	hasAddress     bool
}

func newMessageAssembler() *messageAssembler {
	return &messageAssembler{}
}

// feedCodeword processes a 31-bit BCH-encoded codeword (bit 31 is ignored).
// It BCH-decodes to obtain the 21-bit data word, then processes it:
// Bit 0 = 0: address codeword. Bit 0 = 1: message codeword.
func (ma *messageAssembler) feedCodeword(data uint32) *obj.AlphaMessage {
	// Skip sync and idle codewords
	cw := data & 0x7FFFFFFF
	if cw == syncCodeword || cw == idleCodeword {
		return nil
	}
	// BCH decode to get the 21-bit data word
	decoded, err := bchDecode(data)
	if err != nil {
		return nil
	}
	isMessage := (decoded & 1) == 1
	if isMessage {
		textData := (decoded >> 1) & 0xFFFFF
		ma.accumulateText(textData)
		return nil
	}
	var msg *obj.AlphaMessage
	if ma.hasAddress {
		msg = &obj.AlphaMessage{
			Timestamp: time.Now(),
			CapCode:   ma.currentCapCode,
			Message:   cleanMessage(ma.msgBuf.String()),
			Valid:     ma.msgBuf.Len() > 0,
		}
	}
	capcode := (decoded >> 1) & 0x3FFFF
	funcBits := (decoded >> 19) & 0x3
	ma.currentCapCode = format7Digit(capcode)
	ma.currentFunc = int(funcBits)
	ma.msgBuf.Reset()
	ma.hasAddress = true
	return msg
}

func (ma *messageAssembler) accumulateText(data uint32) {
	text := decode7BitASCII(data, 20)
	ma.msgBuf.WriteString(text)
}

func format7Digit(code uint32) string {
	b := make([]byte, 7)
	for i := 6; i >= 0; i-- {
		b[i] = byte('0' + (code % 10))
		code /= 10
	}
	return string(b)
}

func decode7BitASCII(data uint32, bits int) string {
	var result []byte
	for bits >= 7 {
		ch := byte(data & 0x7F)
		if ch >= 0x20 && ch < 0x7F {
			result = append(result, ch)
		}
		data >>= 7
		bits -= 7
	}
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}

func cleanMessage(msg string) string {
	msg = strings.ReplaceAll(msg, "<NUL>", "")
	msg = strings.ReplaceAll(msg, "<EOT>", "")
	msg = strings.ReplaceAll(msg, "<DC1>", "")
	msg = strings.ReplaceAll(msg, "<DLE>", "")
	msg = strings.ReplaceAll(msg, "<LF>", "|")
	msg = strings.ReplaceAll(msg, "<SUB>J", "|")
	msg = strings.ReplaceAll(msg, "<SUB>M", "|")
	return msg
}

// Test helpers

func buildAddressCodeword(capcode uint32, function int) uint32 {
	data := (capcode << 1) | 0 // bit 0 = 0 (address)
	data |= uint32(function&3) << 19
	return bchEncode(data)
}

func buildMessageCodeword(text string) uint32 {
	var data uint32
	for i := 0; i < len(text) && i < 2; i++ {
		data |= uint32(text[i]&0x7F) << (1 + i*7)
	}
	data |= 1 // bit 0 = 1 (message)
	return bchEncode(data)
}
