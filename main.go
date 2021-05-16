package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/henzai/ranwei/handler"
)

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating Discord session: %v", err)
		return
	}
	defer dg.Close()

	addHandlers(dg)

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
	fmt.Println("再见ranwei")
}

func addHandlers(s *discordgo.Session) {
	// messageCreateEventHandler
	s.AddHandler(handler.OnRecieveAttachments)
	s.AddHandler(handler.ReactionifContainShabu)
	s.AddHandler(handler.ReactionifContainPeyoung)
	// readyEventHandler
	s.AddHandler(handler.OnReady)

	// In this example, we only care about receiving message events.
	s.Identify.Intents = discordgo.IntentsGuildMessages
}
