package internal

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *BotApp) sendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, err := a.bot.Send(msg)
	return err
}

func (a *BotApp) sendMessageWithKeyboard(chatID int64, text string, kb tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = kb
	_, err := a.bot.Send(msg)
	return err
}

func (a *BotApp) answerCallback(query *tgbotapi.CallbackQuery, text string) {
	cb := tgbotapi.NewCallback(query.ID, text)
	if _, err := a.bot.Request(cb); err != nil {
		log.Printf("answer callback: %v", err)
	}
}

func (a *BotApp) editMessageText(chatID int64, messageID int, text string) error {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ParseMode = "Markdown"
	_, err := a.bot.Request(msg)
	return err
}
