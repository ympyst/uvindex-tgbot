package app

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ympyst/uvindex-tgbot/app/handler"
	"github.com/ympyst/uvindex-tgbot/model"
	"github.com/ympyst/uvindex-tgbot/storage"
	"log"
	"os"
)

type App struct {
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

	var s Storage
	var err error

	if os.Getenv("IN_MEMORY_STORAGE") == "1" {
		s = storage.NewMemory()
	} else {
		s, err = storage.NewMongo()
		if err != nil {
			log.Printf("error creating storage: %s", err.Error())
			return nil
		}
	}

	return &App{
		s,
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
	if s.State == model.Start {
		s.Username = a.Telegram.GetUserNameFromUpdate(update)
		s.IsGroup =  a.Telegram.GetIsGroupFromUpdate(update)
	}

	for _, h := range a.handlers {
		h.Handle(ctx, update, s, msg)
	}
	log.Printf( "state after handling update: %v", s)
	err = a.Storage.SaveState(ctx, s)
	if err != nil {
		log.Printf( "error saving new state: %s", err.Error())
	}
}

func (a *App) getState(ctx context.Context, update tgbotapi.Update) (*model.UserState, error) {
	uID := a.Telegram.GetUserIDFromUpdate(update)

	u, err := a.Storage.GetUserStateOrCreate(ctx, uID)
	if err != nil {
		return nil, err
	}

	log.Printf( "state: %v", u)
	return &u, nil
}

