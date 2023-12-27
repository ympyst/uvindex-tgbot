package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	uvindexApp "github.com/ympyst/uvindex-tgbot/app"
	"log"
	"os"
	"strings"
)

const telegramApiTokenEnvKey = "TELEGRAM_API_TOKEN"

type Telegram struct {
	bot *tgbotapi.BotAPI
	app *uvindexApp.App
}

// Button texts
const setLocationBtn = "Set location"
const setAlertsBtn        = "Set alerts"
const setUVIndexThreshold = "Set UV index threshold"

func NewTelegram() *Telegram{
	var err error
	token := os.Getenv(telegramApiTokenEnvKey)
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	if os.Getenv("DEBUG") == "1" {
		bot.Debug = true
	}

	return &Telegram{
		bot: bot,
		app: uvindexApp.NewApp(),
	}
}

func (t *Telegram) Start(ctx context.Context)  {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// `updates` is a golang channel which receives telegram updates
	updates := t.bot.GetUpdatesChan(u)

	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive update from channel and then handle it
		case update := <-updates:
			t.handleUpdate(update)
		}
	}
}

func (t *Telegram) handleUpdate(update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		t.handleMessage(update.Message)
		break

	// Handle button clicks
	case update.CallbackQuery != nil:
		t.handleButton(update.CallbackQuery)
		break
	}
}

func (t *Telegram) handleMessage(message *tgbotapi.Message) {
	user := message.From
	userID := message.From.ID
	text := message.Text
	ctx := context.WithValue(context.Background(), "UserID", userID)

	if user == nil {
		return
	}

	// Print to console
	log.Printf("%s wrote %s", user.FirstName, text)

	var err error
	if strings.HasPrefix(text, "/") {
		err = t.handleCommand(ctx, message.Chat.ID, text)
	} else {
		err = t.app.SetLocation(ctx, text)
	}

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

// When we get a command, we react accordingly
func (t *Telegram) handleCommand(ctx context.Context, chatId int64, command string) error {
	var err error
	var msg tgbotapi.MessageConfig

	switch command {
	case "/uv":
		uv, err := t.app.GetCurrentUVIndex(ctx)
		if err != nil {
			msg = tgbotapi.NewMessage(chatId, err.Error())
		} else {
			msg = tgbotapi.NewMessage(chatId, fmt.Sprintf("%v", uv))
		}
		t.bot.Send(msg)
		break

	case "/menu":
		err = t.sendMenu(chatId)
		break
	}

	return err
}

func (t *Telegram) handleButton(query *tgbotapi.CallbackQuery) {
	switch query.Data {
	case setLocationBtn:
		msg := tgbotapi.NewMessage(query.Message.Chat.ID, "Enter city name")
		t.bot.Send(msg)
		break
	case setAlertsBtn:
		break
	case setUVIndexThreshold:
		break
	}
}


func (t *Telegram) sendMenu(chatId int64) error {
	menuMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(setLocationBtn, setLocationBtn),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(setAlertsBtn, setAlertsBtn),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(setUVIndexThreshold, setUVIndexThreshold),
		),
	)
	msg := tgbotapi.NewMessage(chatId, "Choose setting you want to change")
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = menuMarkup
	_, err := t.bot.Send(msg)
	return err
}
