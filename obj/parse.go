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
	message = strings.ReplaceAll(message, "<NUL>", "")
	message = strings.ReplaceAll(message, "<EOT>", "")
	message = strings.ReplaceAll(message, "<LF>", "|")

	//log.Printf("ts = %d, cap = %s, message = %s", ts.Unix(), cap, message)
	return AlphaMessage{Timestamp: ts, CapCode: cap, Message: message, Valid: message != ""}, nil
}
