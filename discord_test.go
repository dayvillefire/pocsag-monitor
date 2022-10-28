package main

import (
	"testing"
)

func Test_Discord(t *testing.T) {
	return

	_, err := getDiscordClient(discordToken)
	if err != nil {
		t.Fatalf(err.Error())
	}
	_, err = sendDiscordMessage(discordChannel, "Test message")
	if err != nil {
		t.Fatalf(err.Error())
	}
}
