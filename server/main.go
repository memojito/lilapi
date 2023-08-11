package main

import (
	"flag"
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

	address := flag.String("address", "localhost:8080", "The address to listen to")
	flag.Parse()

	session, err := db.NewSession()
	if err != nil {
		log.Printf("Connection failed %v", err)
		return
	}

	api := endpoints.NewAPI(session)

	log.Print("Starting listen on:", *address)
	http.ListenAndServe(*address, api)
}
