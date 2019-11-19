package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

// Player will handle playing sounds to the guild.
// Exits after bot disconnects from the guild.
func (g *Guild) Player() {
	for len(g.Queue) > 0 {
		play := <-g.Queue

		// connect to voice if necessary
		if g.VoiceConnection == nil {
			log.WithFields(log.Fields{
				"guild":   g.Guild.Name,
				"channel": play.Channel.Name,
			}).Debug("Attempting voice connection")
			vc, err := discord.ChannelVoiceJoin(g.Guild.ID, play.Channel.ID, false, false)

			if err != nil {
				log.WithFields(log.Fields{
					"guild":   g.Guild.Name,
					"channel": play.Channel.Name,
					"error":   err,
				}).Error("Voice connection failed")
				return
			}

			log.WithFields(log.Fields{
				"guild":   g.Guild.Name,
				"channel": play.Channel.Name,
			}).Debug("Voice connected")
			g.VoiceConnection = vc
		} else if g.VoiceConnection.ChannelID != play.Channel.ID {
			// change channel if necessary
			log.WithFields(log.Fields{
				"guild":   g.Guild.Name,
				"channel": play.Channel.Name,
			}).Debug("Changing voice channel")
			g.VoiceConnection.ChangeChannel(play.Channel.ID, false, false)
			time.Sleep(time.Millisecond * 125)
		}

		// save stats
		g.SaveToHistory(play)
		go trackSoundStats(play)

		// play sound
		time.Sleep(time.Millisecond * 32)
		play.Skipped = play.Sound.Play(g.VoiceConnection)

		// disconnect if we have forced disconnect pending
		if g.DisconnectPending {
			break
		}

		// enqueue random sound if necessary when state is RNG4EVER
		if g.State == RNG4EVER && len(g.Queue) <= 0 {
			enqueuePlay(nil, play.Guild, play.Sound.Collection, nil)
		}
	}

	log.WithFields(log.Fields{
		"guild": g.Guild.Name,
		"force": g.DisconnectPending,
	}).Info("Disconnecting from voice")
	g.Disconnect()
}
