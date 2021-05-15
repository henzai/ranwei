package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func reactionifContainShabu(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.Contains(m.Content, "シャブ") {
		if ok := shouldReaction(m); !ok {
			return
		}
		err := s.MessageReactionAdd(m.ChannelID, m.ID, "👍")
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot reaction: %v", err)
			return
		}
		return
	}
}

func reactionifContainPeyoung(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.Contains(m.Content, "ペヤング") {
		err := s.MessageReactionAdd(m.ChannelID, m.ID, ":peyoung:")
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot reaction: %v", err)
			return
		}
		return
	}
}

func shouldReaction(m *discordgo.MessageCreate) bool {
	t, err := m.Timestamp.Parse()
	if err != nil {
		return false
	}
	sec := t.Second()
	if res := sec % 2; res != 1 {
		return false
	}
	return true
}
