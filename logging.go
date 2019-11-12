package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	redis "gopkg.in/redis.v3"
)

func trackSoundStats(play *Play) {
	log.WithFields(log.Fields{
		"guild":      play.Guild.Guild.Name,
		"channel":    play.Channel.Name,
		"sound":      play.Sound.Name,
		"collection": play.Sound.Collection.Prefix,
	}).Info("Playing sound")

	if rcli == nil {
		return
	}

	_, err := rcli.Pipelined(func(pipe *redis.Pipeline) error {
		var baseChar string

		if play.Forced {
			baseChar = "specific"
		} else {
			baseChar = "random"
		}

		base := fmt.Sprintf("niksibot:%s", baseChar)
		pipe.Incr("niksibot:total")
		pipe.Incr(fmt.Sprintf("%s:total", base))

		if play.Forced {
			pipe.Incr(fmt.Sprintf("%s:user:%s:sound:%s", base, play.User.ID, play.Sound.Name))
		}

		pipe.Incr(fmt.Sprintf("%s:sound:%s", base, play.Sound.Name))
		pipe.Incr(fmt.Sprintf("%s:guild:%s:sound:%s", base, play.Guild.Guild.ID, play.Sound.Name))
		// pipe.Incr(fmt.Sprintf("%s:guild:%s:chan:%s:sound:%s", base, play.GuildID, play.ChannelID, play.Sound.Name))

		pipe.SAdd(fmt.Sprintf("%s:users", base), play.User.ID)
		pipe.SAdd(fmt.Sprintf("%s:guilds", base), play.Guild.Guild.ID)
		pipe.SAdd(fmt.Sprintf("%s:channels", base), play.Channel.ID)
		return nil
	})

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warning("Failed to track stats in redis")
	}
}

func skipSound(sound *Sound) {
	log.WithFields(log.Fields{
		"sound":      sound.Name,
		"collection": sound.Collection.Prefix,
	}).Debug("Sound skipped")
	if rcli == nil {
		return
	}

	_, err := rcli.Pipelined(func(pipe *redis.Pipeline) error {
		base := "niksibot:skipped"
		pipe.Incr("niksibot:skipped:total")
		// pipe.Incr(fmt.Sprintf("%s:user:%s:sound:%s", base, play.UserID, play.Sound.Name))
		pipe.Incr(fmt.Sprintf("%s:sound:%s", base, sound.Name))
		// pipe.Incr(fmt.Sprintf("%s:guild:%s:sound:%s", base, play.GuildID, play.Sound.Name))
		return nil
	})

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warning("Failed to track stats in redis")
	}

}
