package handler

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	PEYONG_EMOJI_ID = "843157482909728768"
)

func ReactionifContainShabu(s *discordgo.Session, m *discordgo.MessageCreate) {
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

func ReactionifContainPeyoung(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "ペヤング") {
		// NOTE: Privileged Gateway Intents/PRESENCE INTENT がonになっていないとエラーになる
		e, err := s.State.Emoji(m.GuildID, PEYONG_EMOJI_ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot get emoji: %v", err)
			return
		}
		_, err = s.ChannelMessageSend(m.ChannelID, e.MessageFormat())
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot create messege: %v", err)
			return
		}
		// MEMO: reactionを返す時はこのように `name:id`を渡さなくてはいけない
		// see (https://discord.com/developers/docs/resources/channel#create-reaction)
		// err := s.MessageReactionAdd(m.ChannelID, m.ID, "peyoung:843157482909728768")
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "cannot reaction: %v", err)
		// 	return
		// }
		return
	}
}
