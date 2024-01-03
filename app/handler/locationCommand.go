package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ympyst/uvindex-tgbot/model"
)

type LocationCommandHandler struct {

}

func NewLocationCommandHandler() LocationCommandHandler {
	return LocationCommandHandler{}
}

func (h LocationCommandHandler) Handle(ctx context.Context, update tgbotapi.Update, state *model.UserState, msg chan<- tgbotapi.Chattable) {
	if update.Message.Text != "/set_location" {
		return
	}
	state.State = model.WaitingForLocation
	msg <- tgbotapi.NewMessage(update.FromChat().ID, "Please enter your location (e.g. \"Istanbul, Turkey\")")
}
