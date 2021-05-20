package handler

import (
	"errors"
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
		e, err := getPeyongEmoji(s, m.GuildID)
		if err != nil {
			fmt.Printf("cannot get emoji: %v", err)
			return
		}

		_, err = s.ChannelMessageSend(m.ChannelID, e.MessageFormat())
		if err != nil {
			fmt.Printf("cannot create messege: %v", err)
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

func getPeyongEmoji(s *discordgo.Session, guildId string) (*discordgo.Emoji, error) {
	es, err := s.GuildEmojis(guildId)
	if err != nil {
		return nil, err
	}
	for _, ee := range es {
		ee := ee
		if ee.ID == PEYONG_EMOJI_ID {
			return ee, nil
		}
	}
	return nil, errors.New("cannot find peyoung emoji")
}
