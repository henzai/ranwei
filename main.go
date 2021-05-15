package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating Discord session: %v", err)
		return
	}
	defer func() {
		// Cleanly close down the Discord session.
		err = dg.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error cannot close connection: %v", err)
			return
		}
		fmt.Println("再见ranwei")
	}()

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	dg.AddHandler(reactionifContainShabu)
	dg.AddHandler(reactionifContainPeyoung)
	dg.AddHandler(ready)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	ctx := context.Background()
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Attachments) > 0 {
		for _, a := range m.Attachments {
			a := a
			t, err := m.Timestamp.Parse()
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot parse timestamp: %v", err)
				return
			}
			dto := &DTO{
				Attachment: a,
				CreatedAt:  t,
				ChannelID:  m.ChannelID,
				UserID:     m.Author.ID,
				UserName:   m.Author.Username,
			}
			err = saveAttachment(ctx, dto)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot save attachment %v", err)
				return
			}
		}
	}
}

func ready(s *discordgo.Session, r *discordgo.Ready) {
	// Set the playing status.
	s.UpdateGameStatus(0, "シャブスピン")
}
