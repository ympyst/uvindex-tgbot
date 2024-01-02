package app

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ympyst/uvindex-tgbot/model"
	"github.com/ympyst/uvindex-tgbot/weather"
	"log"
	"strings"
)

type LocationHandler struct {
	LocationAPI
}

func NewLocationHandler() LocationHandler {
	return LocationHandler{
		weather.NewLocationAPI(),
	}
}

func (h LocationHandler) Handle(ctx context.Context, update tgbotapi.Update, state *model.UserState, msg chan<- tgbotapi.Chattable)  {
	if update.Message == nil || strings.HasPrefix(update.Message.Text, "/") {
		return
	}
	//if state.State != model.WaitingForLocation {
	//	return
	//}

	ls, err := h.LocationAPI.SearchLocationByName(ctx, update.Message.Text)
	if err != nil {
		log.Printf( "error handling update: %s\n", err.Error())
		return
	}

	var m string

	if len(ls) == 0 {
		m = "No suitable locations found, check spelling or try different location"
	} else if len(ls) > 1 {
		m = "Multiple locations found, try to specify it more precisely (adding country and/or region may help)"
	} else {
		l := ls[0]
		state.Location = l
		state.State = model.Ready
		m = fmt.Sprintf("Location set to: %s, %s, %s", l.Name, l.Region, l.Country)
	}

	log.Printf("state after setting location: %v", state)
	msg <- tgbotapi.NewMessage(update.FromChat().ID, m)
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