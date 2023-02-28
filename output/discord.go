package output

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/dayvillefire/pocsag-monitor/obj"
)

func init() {
	RegisterOutput("discord", func() Output { return &DiscordOutput{} })
}

type DiscordOutput struct {
	discordSession *discordgo.Session
	discordInit    bool
}

func (d *DiscordOutput) Init(token string) error {
	var err error
	if d.discordInit {
		return fmt.Errorf("ERR: already intiialized: %w", err)
	}

	d.discordSession, err = discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("ERR: New(): %w", err)
	}

	err = d.discordSession.Open()
	if err != nil {
		return fmt.Errorf("ERR: Open(): %w", err)
	}

	d.discordInit = true
	return nil
}

func (d *DiscordOutput) SendMessage(a obj.AlphaMessage, channel, msg string) (string, error) {
	// Post normal message
	/*
		res, err := d.discordSession.ChannelMessageSendComplex(channel, &discordgo.MessageSend{
			Content:         msg,
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		})
	*/
	res, err := d.discordSession.ChannelMessageSend(channel, msg)
	if err != nil {
		log.Printf("ERR: SendMessage(): %s", err.Error())
		return "", err
	}

	return res.ID, nil
}
