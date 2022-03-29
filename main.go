package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var stopChan = make(chan bool, 1)

func main() {
	token := os.Getenv("RAT_BOT_TOKEN")
	dg, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		log.Fatal("Error creating Discord session, ", err)
		return
	}

	dg.AddHandler(messageHandler)

	// Only interested in messages
	dg.Identify.Intents = discordgo.IntentGuildMessages | discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuildMembers | discordgo.IntentGuildMembers | discordgo.IntentsAll
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening discord connection, ", err)
		return
	}
	defer dg.Close()

	log.Print("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch m.Content {
	case "!villa":
		PostGigaVilla(s, m)
	case "!råtta", "!rat":
		PostRatImage(s, m)
	case "!gud", "!god":
		PlayLocalAudio(s, m, "./resources/god.dca", stopChan)
	case "!råttparty", "!ratparty", "!råttfest":
		PlayLocalAudio(s, m, "./resources/ratparty.dca", stopChan)
	case "!stop", "!disconnect":
		stopChan <- true
	}

}
