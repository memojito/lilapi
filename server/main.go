package main

import (
	"flag"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"

	"github.com/memojito/lilapi/api/endpoints"
	"github.com/memojito/lilapi/db"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_API_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_API_TOKEN not found")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s LOL", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}

	address := flag.String("address", "localhost:8080", "The address to listen to")
	flag.Parse()

	session, err := db.NewSession()
	if err != nil {
		log.Printf("Connection failed %v", err)
		return
	}

	api := endpoints.NewAPI(session, token)

	log.Print("Starting listen on:", *address)
	http.ListenAndServe(*address, api)
}
