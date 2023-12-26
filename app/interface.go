package app

import (
	"context"
	"github.com/ympyst/uvindex-tgbot/model"
)

type Storage interface {
	GetUserSettingsOrCreate(ctx context.Context, userId int64) (model.UserSettings, error)
	SetUserLocation(ctx context.Context, userId int64, location model.Location) error
}

type WeatherAPI interface {
	GetCurrentUVIndex(ctx context.Context, location model.Location) (float32, error)
}

type LocationAPI interface {
	SearchLocationByName(ctx context.Context, searchQuery string) ([]model.Location, error)
}
