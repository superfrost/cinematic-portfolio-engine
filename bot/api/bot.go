package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"foto-portfolio/bot/internal"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var app = &internal.BotApp{}

func init() {
	app.Init()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("decode update: %v", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := app.HandleUpdate(&update); err != nil {
		log.Printf("handle update: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}
