package main

import (
	"log"
	"os"

	"github.com/memojito/lilapi/bot"
	"github.com/memojito/lilapi/db"
)

func main() {
	url := os.Getenv("POSTGRESQL_URL")
	if url == "" {
		log.Fatal("POSTGRESQL_URL not found")
	}

	conn, err := db.NewConn(url)
	if err != nil {
		log.Printf("Connection failed %v", err)
		return
	}

	token := os.Getenv("TELEGRAM_BOT_API_TOKEN")
	//token := "6387593722:AAEMJFr0OFXKl9Zs1r_-0Peho26_GCQhNCo"
	if token == "" {
		log.Fatal("TELEGRAM_BOT_API_TOKEN not found")
	}

	bot.InitBot(token, conn)
}
