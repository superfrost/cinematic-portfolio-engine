package main

import (
	"log"
	"net/http"
	"os"

	handler "foto-portfolio/bot/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/api/bot", handler.Handler)
	log.Printf("Bot HTTP server listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
