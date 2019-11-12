package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
)

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Content) <= 0 || (m.Content[0] != '!' && len(m.Mentions) < 1) {
		return
	}

	msg := strings.Replace(m.ContentWithMentionsReplaced(), s.State.Ready.User.Username, "username", 1)
	parts := strings.Split(strings.ToLower(msg), " ")

	channel, _ := discord.State.Channel(m.ChannelID)
	if channel == nil {
		log.WithFields(log.Fields{
			"channel": m.ChannelID,
			"message": m.ID,
		}).Warning("Failed to grab channel")
		return
	}

	guild, _ := discord.State.Guild(channel.GuildID)
	if guild == nil {
		log.WithFields(log.Fields{
			"guild":   channel.GuildID,
			"channel": channel,
			"message": m.ID,
		}).Warning("Failed to grab guild")
		return
	}

	// If this is a mention, it should come from the owner (otherwise we don't care)
	if len(m.Mentions) > 0 && m.Author.ID == OWNER && len(parts) > 0 {
		mentioned := false
		for _, mention := range m.Mentions {
			mentioned = (mention.ID == s.State.Ready.User.ID)
			if mentioned {
				break
			}
		}

		if mentioned {
			handleBotControlMessages(s, m, parts, guild)
		}
		return
	}

	if guilds[m.GuildID] == nil {
		guilds[m.GuildID] = &Guild{
			Guild:             guild,
			VoiceConnection:   nil,
			Queue:             nil,
			SkipPending:       false,
			DisconnectPending: false,
			State:             0,
		}
	}
	guildData := guilds[m.GuildID]

	if parts[0] == "!dd" {
		if guildData.VoiceConnection != nil {
			guildData.DisconnectPending = true
			return
		}
	} else if parts[0] == "!skip" {
		guildData.SkipPending = true
	} else if parts[0] == "!history" {
		if guildData.History[0] == nil {
			discord.ChannelMessageSend(channel.ID, "This guild has no history since last bot restart.")
			return
		}

		w := &tabwriter.Writer{}
		buf := &bytes.Buffer{}

		w.Init(buf, 0, 4, 0, ' ', 0)
		fmt.Fprintf(w, ">>> Recently played sounds for current guild:\n")

		for i, el := range guildData.History {
			if el != nil {
				styling := ""
				if el.Skipped {
					styling += "~~"
				}
				if el.Forced {
					styling += "**"
				}

				fmt.Fprintf(w, "%s%d. %s !%s%s\n", styling, i+1, el.Sound.Name, el.Sound.Collection.Prefix, Reverse(styling))
			}
		}

		w.Flush()
		discord.ChannelMessageSend(channel.ID, buf.String())
	}

	// Find the collection for the command we got
	for _, coll := range COLLECTIONS {
		if scontains(parts[0], coll.Commands...) {

			if len(parts) >= 2 && parts[1] == "rng4ever" {
				guildData.State = RNG4EVER
				parts = parts[0:1]
			}

			// If they passed a specific sound effect, find and select that (otherwise play nothing)
			var sound *Sound
			if len(parts) > 1 {
				for _, s := range coll.Sounds {
					if strings.Join(parts[1:len(parts)], " ") == s.Name {
						sound = s
					}
				}

				if sound == nil {
					return
				}
			}

			go enqueuePlay(m.Author, guild, coll, sound)
			return
		}
	}
}

// Reverse the string
// Source: https://stackoverflow.com/questions/1752414/how-to-reverse-a-string-in-go
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
