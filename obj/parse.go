package obj

import (
	"log"
	"regexp"
	"strings"
	"time"
)

var (
	capRegex     *regexp.Regexp
	messageRegex *regexp.Regexp
)

func init() {
	var err error
	capRegex, err = regexp.Compile(` Address:\s+(\d+) `)
	if err != nil {
		log.Fatal(err)
	}
	messageRegex, err = regexp.Compile(` Alpha:   (.*)`)
	if err != nil {
		log.Fatal(err)
	}
}

func ParseAlphaMessage(ts time.Time, m string) (AlphaMessage, error) {
	if !strings.Contains(m, "POCSAG512: ") {
		log.Printf("DEBUG: %s\n", m)
		return AlphaMessage{}, nil
	}
	//log.Printf("INFO: %s\n", m)

	// Parse cap and message
	var cap string
	capMatches := capRegex.FindStringSubmatch(m)
	if len(capMatches) > 1 {
		cap = capMatches[1]
	}
	var message string
	messageMatches := messageRegex.FindStringSubmatch(m)
	if len(messageMatches) > 1 {
		message = messageMatches[1]
	}

	// Remove all null, etc
	message = runMultiple(message, strings.ReplaceAll, map[string]string{
		"<NUL>":  "",
		"<EOT>":  "",
		"<DC1>":  "",
		"<DLE>":  "",
		"<LF>":   "|",
		"<SUB>J": "|",
		"<SUB>M": "|",
	})
	/*
		message = strings.ReplaceAll(message, "<NUL>", "")
		message = strings.ReplaceAll(message, "<EOT>", "")
		message = strings.ReplaceAll(message, "<LF>", "|")
		message = strings.ReplaceAll(message, "<SUB>J", "|")
		message = strings.ReplaceAll(message, "<SUB>M", "|")
	*/

	//log.Printf("ts = %d, cap = %s, message = %s", ts.Unix(), cap, message)
	return AlphaMessage{Timestamp: ts, CapCode: cap, Message: message, Valid: message != ""}, nil
}

func runMultiple(orig string, fn func(string, string, string) string, data map[string]string) string {
	if fn == nil {
		return orig
	}
	x := orig
	for k, v := range data {
		x = fn(x, k, v)
	}
	return x
}
 