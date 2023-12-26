package main

import (
	"context"
	"fmt"
	uvindexApp "github.com/ympyst/uvindex-tgbot/app"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	a *uvindexApp.App
	// Button texts
	setLocationBtn = "Set location"
	setAlertsBtn        = "Set alerts"
	setUVIndexThreshold = "Set UV index threshold"

	// Store current uv index
	uv float32
	bot       *tgbotapi.BotAPI

	// Keyboard layout for the first menu. One button, one row
	menuMarkup = tgbotapi.NewInlineKeyboardMarkup(
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
)

func main() {
	var err error
	token := os.Getenv("TELEGRAM_API_TOKEN")
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}

	a = uvindexApp.NewApp()

	// Set this to true to log all interactions with telegram servers
	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	ctx := context.Background()

	// `updates` is a golang channel which receives telegram updates
	updates := bot.GetUpdatesChan(u)

	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive update from channel and then handle it
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		handleMessage(update.Message)
		break

	// Handle button clicks
	case update.CallbackQuery != nil:
		handleButton(update.CallbackQuery)
		break
	}
}

func handleMessage(message *tgbotapi.Message) {
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
		err = handleCommand(ctx, message.Chat.ID, text)
	} else {
		err = a.SetLocation(ctx, text)
	}

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

// When we get a command, we react accordingly
func handleCommand(ctx context.Context, chatId int64, command string) error {
	var err error
	var msg tgbotapi.MessageConfig

	switch command {
	case "/uv":
		uv, err = a.GetCurrentUVIndex(ctx)
		if err != nil {
			msg = tgbotapi.NewMessage(chatId, err.Error())
		} else {
			msg = tgbotapi.NewMessage(chatId, fmt.Sprintf("%v", uv))
		}
		bot.Send(msg)
		break

	case "/menu":
		err = sendMenu(chatId)
		break
	}

	return err
}

func handleButton(query *tgbotapi.CallbackQuery) {
	switch query.Data {
	case setLocationBtn:
		msg := tgbotapi.NewMessage(query.Message.Chat.ID, "Enter city name")
		bot.Send(msg)
		break
	case setAlertsBtn:
		break
	case setUVIndexThreshold:
		break
	}
}

func sendMenu(chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, "Choose setting you want to change")
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = menuMarkup
	_, err := bot.Send(msg)
	return err
}

