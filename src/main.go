package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cox/src/handlers"
	"cox/src/utils"

	"github.com/bwmarrin/discordgo"
)

func main() {

	token := utils.GetEnv("DISCORD_BOT_TOKEN", true)

	dg, err := discordgo.New(token)

	if err != nil {
		log.Println(err)
		return
	}

	// Add a handler for new messages
	dg.AddHandler(handlers.MessageCreateHandler)
	dg.AddHandler(handlers.MessageReactionAddHandler)

	// Open a websocket connection to Discord
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer dg.Close()

	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	// Wait here until CTRL-C or other term signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

}
