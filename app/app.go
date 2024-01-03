package app

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ympyst/uvindex-tgbot/app/handler"
	"github.com/ympyst/uvindex-tgbot/model"
	"github.com/ympyst/uvindex-tgbot/storage"
	"github.com/ympyst/uvindex-tgbot/weather"
	"log"
)

const UserIDCtxKey = "UserID"

type App struct {
	WeatherAPI
	LocationAPI
	Storage
	*Telegram
	handlers []UpdateHandler
}

func NewApp() *App {
	h := []UpdateHandler{
		handler.NewStartHandler(),
		handler.NewLocationHandler(),
		handler.NewCurrentUVIndexHandler(),
		handler.NewLocationCommandHandler(),
	}

	return &App{
		weather.NewAPI(),
		weather.NewLocationAPI(),
		storage.NewStorage(),
		NewTelegram(),
		h,
	}
}

func (a *App) Start(ctx context.Context)  {
	updates := a.Telegram.GetUpdatesChan()
	messages := make(chan tgbotapi.Chattable, 10)

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			a.handleUpdate(ctx, update, messages)
		case msg := <-messages:
			_, err := a.Telegram.bot.Send(msg)
			if err != nil {
				log.Printf("error sending message: %s", err.Error())
			}
		}
	}
}

func (a *App) handleUpdate(ctx context.Context, update tgbotapi.Update, msg chan<- tgbotapi.Chattable)  {
	s, err := a.getState(ctx, update)
	if err != nil {
		log.Printf( "error handling update: %s", err.Error())
	}

	for _, h := range a.handlers {
		h.Handle(ctx, update, &s, msg)
	}
	log.Printf( "state after handling update: %v", s)
	err = a.Storage.SaveState(ctx, &s)
	if err != nil {
		log.Printf( "error saving new state: %s", err.Error())
	}
}

func (a *App) getState(ctx context.Context, update tgbotapi.Update) (model.UserState, error) {
	uID := a.Telegram.GetUserIDFromUpdate(update)

	u, err := a.Storage.GetUserSettingsOrCreate(ctx, uID)
	if err != nil {
		return model.UserState{}, err
	}

	log.Printf( "state: %v", u)
	return u, nil
}

