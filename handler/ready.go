package handler

import "github.com/bwmarrin/discordgo"

func OnReady(s *discordgo.Session, r *discordgo.Ready) {
	s.UpdateGameStatus(0, "シャブスピン")
}
