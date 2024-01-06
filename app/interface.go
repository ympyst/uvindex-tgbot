package app

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ympyst/uvindex-tgbot/model"
)

type Storage interface {
	GetUserSettingsOrCreate(ctx context.Context, userId int64) (model.UserState, error)
	SaveState(ctx context.Context, state *model.UserState) error
}

type UpdateHandler interface {
	Handle(ctx context.Context, update tgbotapi.Update, state *model.UserState, msg chan<- tgbotapi.Chattable)
}
