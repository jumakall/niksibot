package main

import "github.com/bwmarrin/discordgo"

// Guild wraps Discord's guild and adds extra info that should carry over the app
type Guild struct {
	Guild             *discordgo.Guild
	VoiceConnection   *discordgo.VoiceConnection
	Queue             chan *Play
	History           [MAX_HISTORY_SIZE]*Play
	SkipPending       bool
	DisconnectPending bool
	State             int
}

// Disconnect guild's voice connection
func (g *Guild) Disconnect() {
	g.VoiceConnection.Disconnect()
	g.VoiceConnection = nil
	g.Reset()
}

// Reset guild's queue and pending operations
// Called when the bot disconnects
func (g *Guild) Reset() {
	g.Queue = nil
	g.DisconnectPending = false
	g.SkipPending = false
	g.State = 0
}

// SaveToHistory adds the play to the guild's play history
// Automatically removes oldest play when full
func (g *Guild) SaveToHistory(p *Play) {
	for i := len(g.History) - 1; i > 0; i-- {
		g.History[i] = g.History[i-1]
	}
	g.History[0] = p
}
