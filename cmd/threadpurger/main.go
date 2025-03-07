package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	channelId   = flag.String("channel", "1031981006883913779", "Channel ID")
	startMsg    = flag.String("start", "", "Start message ID")
	matchText   = flag.String("matchText", "187001: ", "Text to match to delete messages. Leave blank to delete all.")
	threadPurge = flag.Bool("threadPurge", false, "Also purge threads explicitly")
	dryrun      = flag.Bool("dryrun", false, "Don't actually do")

	discordSession *discordgo.Session
	discordInit    bool
)

func main() {
	flag.Parse()

	t := true

	err := initDiscord(token)
	if err != nil {
		panic(err)
	}

	beforeId := *startMsg
	for {
		msgs, err := discordSession.ChannelMessages(*channelId, 100, beforeId, "", "")
		if err != nil {
			panic(err)
		}

		log.Printf("Found %d channel messages from point '%s' (max batch 100)", len(msgs), beforeId)

		if len(msgs) < 1 {
			break
		}

		del := make([]string, 0)

		for k, v := range msgs {
			if v.Thread != nil {
				if v.Thread.MessageCount > 2 {
					log.Printf("Starts thread, skipping")
					continue
				}
			}

			beforeId = v.ID

			if *matchText != "" {
				found := false
				if strings.Contains(*matchText, "|") {
					sub := strings.Split(*matchText, "|")
					for _, x := range sub {
						if strings.Contains(v.Content, x) {
							found = true
						}
					}
				} else {
					if !strings.Contains(v.Content, *matchText) {
						continue
					}
				}
				if !found {
					continue
				}
			}

			log.Printf("Found message #%d, text %s, ID %s", k, v.Content, v.ID)
			// Add to list of messages to bulk delete
			del = append(del, v.ID)
		}

		if len(del) == 0 {
			//break
		}

		if !*dryrun {
			log.Printf("Purging messages #%v", del)
			for _, msg := range del {
				err = discordSession.ChannelMessageDelete(*channelId, msg)
				if err != nil {
					panic(err)
				}
			}
			//err = discordSession.ChannelMessagesBulkDelete(*channelId, del)
		}
	}

	if !*threadPurge {
		log.Printf("Set not to purge threads explicitly, exiting")
		return
	}

	for {
		threads, err := discordSession.ThreadsActive(*channelId)
		if err != nil {
			panic(err)
		}

		log.Printf("Threads: %d found", len(threads.Threads))
		for k, v := range threads.Threads {
			log.Printf("[%d] %s (%d messages)", k, v.Topic, v.MessageCount)
			if v.MessageCount > 2 {
				/*
					if v.ThreadMetadata.AutoArchiveDuration != 300 {
						_, err := discordSession.ChannelEdit(v.ID, &discordgo.ChannelEdit{AutoArchiveDuration: 300})
						if err != nil {
							log.Printf("ERR: %s", err.Error())
						}
					}
				*/
				continue
			}
			_, err := discordSession.ChannelEdit(v.ID, &discordgo.ChannelEdit{Archived: &t})
			if err != nil {
				log.Printf("ERR: %s", err.Error())
			}
		}

		if !threads.HasMore {
			log.Printf("No more threads, exiting")
			break
		}

		time.Sleep(time.Second)

		log.Printf("More threads, continuing")
	}

}

func initDiscord(token string) error {
	var err error
	if discordInit {
		return fmt.Errorf("ERR: already intiialized: %w", err)
	}

	discordSession, err = discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("ERR: discordgo.New(): %w", err)
	}

	err = discordSession.Open()
	if err != nil {
		return fmt.Errorf("ERR: discordSession.Open(): %w", err)
	}

	discordInit = true
	return nil
}

func derefChannelMap(in map[int64]*discordgo.Channel) map[int64]discordgo.Channel {
	out := map[int64]discordgo.Channel{}
	for k, v := range in {
		out[k] = *v
	}
	return out
}

func rerefChannelMap(in map[int64]discordgo.Channel) map[int64]*discordgo.Channel {
	out := map[int64]*discordgo.Channel{}
	for k, v := range in {
		out[k] = &v
	}
	return out
}
