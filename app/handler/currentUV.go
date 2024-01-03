package handler

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ympyst/uvindex-tgbot/model"
	"github.com/ympyst/uvindex-tgbot/weather"
	"log"
)

type WeatherAPI interface {
	GetCurrentUVIndex(ctx context.Context, location model.Location) (float32, error)
}

type CurrentUVIndexHandler struct {
	WeatherAPI
}

func NewCurrentUVIndexHandler() CurrentUVIndexHandler {
	return CurrentUVIndexHandler{
		weather.NewAPI(),
	}
}

func (h CurrentUVIndexHandler) Handle(ctx context.Context, update tgbotapi.Update, state *model.UserState, msg chan<- tgbotapi.Chattable)  {
	if state.State != model.Ready || update.Message == nil || update.Message.Text != "/uv" {
		return
	}
	uv, err := h.WeatherAPI.GetCurrentUVIndex(ctx, state.Location)
	if err != nil {
		log.Printf("error getting current UV index: %s", err.Error())
	}

	msg <- tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("%v", uv))
}