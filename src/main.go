package main

import (
	"context"
	"cox/src/handlers"
	"cox/src/utils"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func discordBotSession(ctx context.Context, wg *sync.WaitGroup, subChannel <-chan *discordgo.MessageCreate) {
	defer wg.Done()

	fmt.Printf("Bot Session Started\n")
	token := utils.GetEnv("DISCORD_BOT_TOKEN", true)

	dg, err := discordgo.New(token)

	if err != nil {
		log.Println(err)
		return
	}

	// Add a handler for new messages
	dg.AddHandler(handlers.MessageCreateHandler)
	dg.AddHandler(handlers.MessageReactionAddHandler)
	dg.AddHandler(handlers.MessageReactionRemoveHandler)

	// Add handler from messages received by other account
	handlers.BufferedMessageCreateHandler(subChannel, dg)

	// Open a websocket connection to Discord
	err = dg.Open()

	if err != nil {
		fmt.Println("Error establishing connection with discord.", err)
		return
	}

	defer dg.Close()

	<-ctx.Done()
	fmt.Printf("Bot Session Closed\n")
}

func accountSession(ctx context.Context, wg *sync.WaitGroup, pubChannel chan<- *discordgo.MessageCreate) {
	defer wg.Done()

	fmt.Printf("Account Session Started\n")
	token := utils.GetEnv("DISCORD_ACCOUNT_TOKEN", true)
	dg, err := discordgo.New(token)
	dg.Identify.Intents = discordgo.IntentGuildMessages | discordgo.IntentMessageContent

	if err != nil {
		log.Println(err)
		return
	}

	// Add a handler for new messages
	dg.AddHandler(handlers.AccountMessageCreateHandler(pubChannel))

	// Open a websocket connection to Discord
	err = dg.Open()

	if err != nil {
		fmt.Println("Error establishing connection with discord.", err)
		return
	}

	defer dg.Close()

	<-ctx.Done()
	fmt.Printf("Account Session Closed\n")
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	// Wait here until CTRL-C or other term signal
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	msgChannel := make(chan *discordgo.MessageCreate, 100)

	go discordBotSession(ctx, &wg, msgChannel)
	go accountSession(ctx, &wg, msgChannel)

	// Wait for cancellation
	<-ctx.Done()
	fmt.Println("Shutting down...")

	// Wait for both goroutines to finish
	wg.Wait()
	fmt.Println("All workers stopped.")
}
