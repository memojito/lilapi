package main

import (
	"flag"
	"github.com/memojito/lilapi/api/endpoints"
	"github.com/memojito/lilapi/db"
	"log"
	"net/http"
)

func main() {
	address := flag.String("address", "localhost:8080", "The address to listen to")

	session, err := db.NewSession()
	if err != nil {
		log.Printf("Connection failed")
	}

	api := endpoints.NewAPI(session)

	http.ListenAndServe(*address, api)
}
