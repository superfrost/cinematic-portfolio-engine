package internal

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *BotApp) cmdStart(msg *tgbotapi.Message) error {
	text := `*Portfolio Content Manager Bot*

Manage your portfolio content via Telegram!

*Commands:*
/cases — List all cases
/reviews — List all reviews
/add — Add new case or review
/status — Last deploy status
/start — Show this help`
	return a.sendMessage(msg.Chat.ID, text)
}

func (a *BotApp) cmdListCases(msg *tgbotapi.Message) error {
	return a.listContent(msg, "src/lib/content/cases/")
}

func (a *BotApp) cmdListReviews(msg *tgbotapi.Message) error {
	return a.listContent(msg, "src/lib/content/reviews/")
}

func (a *BotApp) listContent(msg *tgbotapi.Message, dirPath string) error {
	files, err := a.gh.ListFiles(a.ctx, dirPath)
	if err != nil {
		return a.sendMessage(msg.Chat.ID, fmt.Sprintf("Error listing files: %v", err))
	}

	if len(files) == 0 {
		return a.sendMessage(msg.Chat.ID, "No files found.")
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, f := range files {
		name := strings.TrimSuffix(f.Name, ".json")
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📄 "+name, "noop_"+f.Path),
			tgbotapi.NewInlineKeyboardButtonData("Edit", "edit_"+f.Path),
			tgbotapi.NewInlineKeyboardButtonData("Rename", "rename_"+f.Path),
			tgbotapi.NewInlineKeyboardButtonData("Delete", "delete_"+f.Path),
		)
		rows = append(rows, row)
	}

	kb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return a.sendMessageWithKeyboard(msg.Chat.ID, fmt.Sprintf("*%s*", dirPath), kb)
}

func (a *BotApp) cmdAdd(msg *tgbotapi.Message) error {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📁 Case", "add_case"),
			tgbotapi.NewInlineKeyboardButtonData("💬 Review", "add_review"),
		),
	)
	return a.sendMessageWithKeyboard(msg.Chat.ID, "What would you like to add?", kb)
}

func (a *BotApp) cmdStatus(msg *tgbotapi.Message) error {
	ds, err := a.gh.GetDeployStatus(a.ctx)
	if err != nil {
		return a.sendMessage(msg.Chat.ID, fmt.Sprintf("Error fetching deploy status: %v", err))
	}

	statusIcon := "❌"
	if ds.Success {
		statusIcon = "✅"
	}

	text := fmt.Sprintf("*Deploy Status*\n%s Conclusion: *%s*\nUpdated: `%s`",
		statusIcon, ds.Conclusion, ds.UpdatedAt.Format("2006-01-02 15:04:05"))
	return a.sendMessage(msg.Chat.ID, text)
}

func (a *BotApp) createCaseFromState(chatID int64, userID int64, us *UserState) error {
	slug := Transliterate(us.Data["title"])
	date := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s-%s.json", date, slug)
	path := fmt.Sprintf("src/lib/content/cases/%s", filename)
	id := strings.TrimSuffix(filename, ".json")

	caseData := map[string]interface{}{
		"id":          id,
		"title":       us.Data["title"],
		"description": us.Data["description"],
		"items": []map[string]interface{}{
			{
				"type":     "video",
				"provider": us.Data["provider"],
				"id":       us.Data["video_id"],
			},
		},
	}

	content, _ := json.MarshalIndent(caseData, "", "    ")
	if err := a.gh.CreateFile(a.ctx, path, string(content), fmt.Sprintf("Add case %s via bot", filename)); err != nil {
		a.clearUserState(userID)
		return a.sendMessage(chatID, fmt.Sprintf("Error creating file: %v", err))
	}

	a.clearUserState(userID)
	return a.sendMessage(chatID, fmt.Sprintf("Case created: `%s`", path))
}

func (a *BotApp) createReviewFromState(chatID int64, userID int64, us *UserState) error {
	slug := Transliterate(us.Data["author"])
	date := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s-%s.json", date, slug)
	path := fmt.Sprintf("src/lib/content/reviews/%s", filename)

	reviewData := map[string]interface{}{
		"author": us.Data["author"],
		"text":   us.Data["text"],
	}

	if us.Data["provider"] != "" {
		reviewData["video"] = map[string]string{
			"provider": us.Data["provider"],
			"id":       us.Data["video_id"],
		}
	}

	content, _ := json.MarshalIndent(reviewData, "", "    ")
	if err := a.gh.CreateFile(a.ctx, path, string(content), fmt.Sprintf("Add review %s via bot", filename)); err != nil {
		a.clearUserState(userID)
		return a.sendMessage(chatID, fmt.Sprintf("Error creating file: %v", err))
	}

	a.clearUserState(userID)
	return a.sendMessage(chatID, fmt.Sprintf("Review created: `%s`", path))
}

func (a *BotApp) handleCallback(query *tgbotapi.CallbackQuery) error {
	data := query.Data
	chatID := query.Message.Chat.ID
	userID := query.From.ID
	msgID := query.Message.MessageID

	switch {
	case data == "add_case":
		us := a.getUserState(userID)
		us.State = StateAwaitingCaseLink
		us.Data = make(map[string]string)
		a.answerCallback(query, "Send a YouTube or Rutube link")
		return a.editMessageText(chatID, msgID, "Send a video link (YouTube or Rutube) for the case:")

	case data == "add_review":
		us := a.getUserState(userID)
		us.State = StateAwaitingReviewLink
		us.Data = make(map[string]string)
		a.answerCallback(query, "Send a video link or /skip")
		return a.editMessageText(chatID, msgID, "Send a video link (YouTube or Rutube) for the review, or send /skip to skip video:")

	case data == "cancel_delete":
		a.answerCallback(query, "Cancelled")
		return a.editMessageText(chatID, msgID, "Deletion cancelled.")

	case strings.HasPrefix(data, "delete_"):
		path := strings.TrimPrefix(data, "delete_")
		kb := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅ Yes, delete", "confirm_delete_"+path),
				tgbotapi.NewInlineKeyboardButtonData("❌ Cancel", "cancel_delete"),
			),
		)
		a.answerCallback(query, "Confirm deletion")
		msg := tgbotapi.NewEditMessageText(chatID, msgID, fmt.Sprintf("Are you sure you want to delete `%s`?", path))
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = &kb
		if _, err := a.bot.Request(msg); err != nil {
			return err
		}
		return nil

	case strings.HasPrefix(data, "confirm_delete_"):
		path := strings.TrimPrefix(data, "confirm_delete_")
		fc, err := a.gh.GetFile(a.ctx, path)
		if err != nil {
			a.answerCallback(query, "Error reading file")
			return a.editMessageText(chatID, msgID, fmt.Sprintf("Error reading file: %v", err))
		}
		if err := a.gh.DeleteFile(a.ctx, path, fc.SHA, fmt.Sprintf("Delete %s via bot", path)); err != nil {
			a.answerCallback(query, "Error deleting")
			return a.editMessageText(chatID, msgID, fmt.Sprintf("Error deleting: %v", err))
		}
		a.answerCallback(query, "Deleted")
		return a.editMessageText(chatID, msgID, fmt.Sprintf("Deleted `%s`", path))

	case strings.HasPrefix(data, "edit_"):
		path := strings.TrimPrefix(data, "edit_")
		us := a.getUserState(userID)
		us.State = StateAwaitingEdit
		us.Data = map[string]string{"file_path": path}
		a.answerCallback(query, "Send new content")
		return a.editMessageText(chatID, msgID, "Send the new content (description for case, text for review):")

	case strings.HasPrefix(data, "rename_"):
		path := strings.TrimPrefix(data, "rename_")
		us := a.getUserState(userID)
		us.State = StateAwaitingRename
		us.Data = map[string]string{"file_path": path}
		a.answerCallback(query, "Send new name")
		return a.editMessageText(chatID, msgID, "Send the new name (will be transliterated, date preserved):")

	case strings.HasPrefix(data, "noop_"):
		a.answerCallback(query, "")
		return nil
	}

	return nil
}
