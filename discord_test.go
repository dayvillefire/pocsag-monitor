package main

import "testing"

func Test_Discord(t *testing.T) {
	_, err := getDiscordClient(discordToken)
	if err != nil {
		t.Fatalf(err.Error())
	}
	_, err = sendDiscordMessage("Test message")
	if err != nil {
		t.Fatalf(err.Error())
	}
}
