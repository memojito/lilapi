package main

import (
	"github.com/memojito/lilapi/bot"
	"github.com/memojito/lilapi/db"
	"log"
	"os"
)

func main() {
	session, err := db.NewSession()
	if err != nil {
		log.Printf("Connection failed %v", err)
		return
	}

	token := os.Getenv("TELEGRAM_BOT_API_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_API_TOKEN not found")
	}

	bot.InitBot(token, session)
}
