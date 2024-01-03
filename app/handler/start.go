package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ympyst/uvindex-tgbot/model"
)

type StartHandler struct {

}

func NewStartHandler() StartHandler {
	return StartHandler{}
}

func (h StartHandler) Handle(ctx context.Context, update tgbotapi.Update, state *model.UserState, msg chan<- tgbotapi.Chattable) {
	if update.Message.Text != "/start" {
		return
	}
	state.State = model.WaitingForLocation
	msg <- tgbotapi.NewMessage(update.FromChat().ID, "Please enter your location (e.g. \"Istanbul, Turkey\")")
}
