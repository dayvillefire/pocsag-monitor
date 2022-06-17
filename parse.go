package main

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

type alphaMessage struct {
	Timestamp time.Time
	CapCode   string
	Message   string
	Valid     bool
}

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

func parse(ts time.Time, m string) (alphaMessage, error) {
	if !strings.Contains(m, "POCSAG512: ") {
		log.Printf("%s: DEBUG: %s\n", ts.Format("2006-01-02 15:04:05"), m)
		return alphaMessage{}, nil
	}
	//log.Printf("%s: INFO: %s\n", ts.Format("2006-01-02 15:04:05"), m)

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
	return alphaMessage{Timestamp: ts, CapCode: cap, Message: message, Valid: message != ""}, nil
}
