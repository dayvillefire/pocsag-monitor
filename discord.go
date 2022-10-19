package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	discordSession *discordgo.Session
	discordInit    bool
)

func getDiscordClient(token string) (*discordgo.Session, error) {
	var err error
	if discordInit {
		return discordSession, fmt.Errorf("ERR: already intiialized: %w", err)
	}

	discordSession, err = discordgo.New("Bot " + token)
	if err != nil {
		return discordSession, fmt.Errorf("ERR: New(): %w", err)
	}

	err = discordSession.Open()
	if err != nil {
		return discordSession, fmt.Errorf("ERR: Open(): %w", err)
	}

	/*
		userinfo, err := discordSession.User("@me")
		if err != nil {
			return discordSession, fmt.Errorf("ERR: User(): %w", err)
		}
		log.Printf("%#v", userinfo)
	*/

	discordInit = true
	return discordSession, err
}

func sendDiscordMessage(msg string) (string, error) {
	m := discordgo.MessageSend{
		Content:         msg,
		AllowedMentions: &discordgo.MessageAllowedMentions{},
	}

	// Post normal message
	res, err := discordSession.ChannelMessageSendComplex(discordChannel, &m)
	if err != nil {
		return "", err
	}

	return res.ID, nil
}
