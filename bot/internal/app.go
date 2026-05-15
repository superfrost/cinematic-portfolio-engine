package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type BotApp struct {
	bot    *tgbotapi.BotAPI
	gh     *GitHubClient
	ctx    context.Context
	owners map[int64]bool
	mu     sync.Mutex
	states map[int64]*UserState
}

func (a *BotApp) Init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN is required")
	}
	var err error
	a.bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("NewBotAPI: %v", err)
	}

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("GITHUB_TOKEN is required")
	}

	ghOwner := os.Getenv("GITHUB_OWNER")
	if ghOwner == "" {
		ghOwner = "superfrost"
	}
	ghRepo := os.Getenv("GITHUB_REPO")
	if ghRepo == "" {
		ghRepo = "cinematic-portfolio-engine"
	}

	a.ctx = context.Background()
	a.gh = NewGitHubClient(a.ctx, githubToken, ghOwner, ghRepo)

	a.owners = make(map[int64]bool)
	allowed := os.Getenv("ALLOWED_USER_ID")
	for _, idStr := range strings.Split(allowed, ",") {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		var id int64
		if _, err := fmt.Sscanf(idStr, "%d", &id); err == nil {
			a.owners[id] = true
		}
	}

	a.states = make(map[int64]*UserState)
	a.bot.Debug = false
}

func (a *BotApp) isAllowed(userID int64) bool {
	if len(a.owners) == 0 {
		return true
	}
	return a.owners[userID]
}

func (a *BotApp) getUserState(userID int64) *UserState {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.states[userID]; !ok {
		a.states[userID] = &UserState{Data: make(map[string]string)}
	}
	return a.states[userID]
}

func (a *BotApp) clearUserState(userID int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.states, userID)
}

func (a *BotApp) HandleUpdate(update *tgbotapi.Update) error {
	if update.Message != nil {
		if !a.isAllowed(update.Message.From.ID) {
			log.Printf("unauthorized message: %d", update.Message.From.ID)
			return nil
		}
		return a.handleMessage(update.Message)
	}
	if update.CallbackQuery != nil {
		if !a.isAllowed(update.CallbackQuery.From.ID) {
			log.Printf("unauthorized callback: %d", update.CallbackQuery.From.ID)
			return nil
		}
		return a.handleCallback(update.CallbackQuery)
	}
	return nil
}

func (a *BotApp) handleMessage(msg *tgbotapi.Message) error {
	userID := msg.From.ID
	us := a.getUserState(userID)

	if us.State != StateNone {
		return a.handleStatefulMessage(msg, us)
	}

	if msg.IsCommand() {
		switch msg.Command() {
		case "start":
			return a.cmdStart(msg)
		case "cases":
			return a.cmdListCases(msg)
		case "reviews":
			return a.cmdListReviews(msg)
		case "add":
			return a.cmdAdd(msg)
		case "status":
			return a.cmdStatus(msg)
		}
	}

	return a.sendMessage(msg.Chat.ID, "Unknown command. Use /start for help.")
}

func (a *BotApp) handleStatefulMessage(msg *tgbotapi.Message, us *UserState) error {
	chatID := msg.Chat.ID
	userID := msg.From.ID

	switch us.State {
	case StateAwaitingCaseLink:
		video, err := ParseVideoLink(msg.Text)
		if err != nil {
			return a.sendMessage(chatID, "Invalid video link. Send a YouTube or Rutube URL.")
		}
		us.Data["provider"] = video.Provider
		us.Data["video_id"] = video.ID
		us.State = StateAwaitingCaseTitle
		return a.sendMessage(chatID, "Link saved! Now send the case title:")

	case StateAwaitingCaseTitle:
		us.Data["title"] = msg.Text
		us.State = StateAwaitingCaseDescription
		return a.sendMessage(chatID, "Title saved! Now send the case description:")

	case StateAwaitingCaseDescription:
		us.Data["description"] = msg.Text
		return a.createCaseFromState(chatID, userID, us)

	case StateAwaitingReviewLink:
		if msg.Text == "/skip" {
			us.Data["provider"] = ""
			us.Data["video_id"] = ""
		} else {
			video, err := ParseVideoLink(msg.Text)
			if err != nil {
				return a.sendMessage(chatID, "Invalid video link. Send a YouTube or Rutube URL, or /skip to skip video.")
			}
			us.Data["provider"] = video.Provider
			us.Data["video_id"] = video.ID
		}
		us.State = StateAwaitingReviewAuthor
		return a.sendMessage(chatID, "Now send the client/author name:")

	case StateAwaitingReviewAuthor:
		us.Data["author"] = msg.Text
		us.State = StateAwaitingReviewText
		return a.sendMessage(chatID, "Author saved! Now send the review text:")

	case StateAwaitingReviewText:
		us.Data["text"] = msg.Text
		return a.createReviewFromState(chatID, userID, us)

	case StateAwaitingEdit:
		path := us.Data["file_path"]
		fc, err := a.gh.GetFile(a.ctx, path)
		if err != nil {
			a.clearUserState(userID)
			return a.sendMessage(chatID, fmt.Sprintf("Error reading file: %v", err))
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(fc.Content), &data); err != nil {
			a.clearUserState(userID)
			return a.sendMessage(chatID, "Error parsing JSON")
		}

		if strings.Contains(path, "cases") {
			data["description"] = msg.Text
		} else {
			data["text"] = msg.Text
		}

		newContent, _ := json.MarshalIndent(data, "", "    ")
		if err := a.gh.UpdateFile(a.ctx, path, fc.SHA, string(newContent), fmt.Sprintf("Update %s via bot", path)); err != nil {
			a.clearUserState(userID)
			return a.sendMessage(chatID, fmt.Sprintf("Error updating: %v", err))
		}

		a.clearUserState(userID)
		return a.sendMessage(chatID, "File updated successfully!")

	case StateAwaitingRename:
		oldPath := us.Data["file_path"]
		parsed, ok := ParseFilename(oldPath)
		if !ok {
			a.clearUserState(userID)
			return a.sendMessage(chatID, "Invalid filename format")
		}

		newSlug := Transliterate(msg.Text)
		newPath := strings.Replace(oldPath, parsed.Date+"-"+parsed.Slug, parsed.Date+"-"+newSlug, 1)

		if err := a.gh.RenameFile(a.ctx, oldPath, newPath, fmt.Sprintf("Rename %s -> %s via bot", oldPath, newPath)); err != nil {
			a.clearUserState(userID)
			return a.sendMessage(chatID, fmt.Sprintf("Error renaming: %v", err))
		}

		a.clearUserState(userID)
		return a.sendMessage(chatID, fmt.Sprintf("Renamed to `%s`!", newPath))
	}

	return nil
}
