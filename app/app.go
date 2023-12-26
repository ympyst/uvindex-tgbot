package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/ympyst/uvindex-tgbot/storage"
	"github.com/ympyst/uvindex-tgbot/weather"
)

const UserIDCtxKey = "UserID"

type App struct {
	WeatherAPI
	LocationAPI
	Storage
}

func NewApp() *App {
	return &App{
		weather.NewAPI(),
		weather.NewLocationAPI(),
		storage.NewStorage(),
	}
}

func (a *App) SetLocation(ctx context.Context, searchQuery string) error {
	uID, err := a.getUserIDFromCtx(ctx)
	if err != nil {
		return err
	}

	l, err := a.LocationAPI.SearchLocationByName(ctx, searchQuery)
	if err != nil {
		return err
	}
	if len(l) == 0 {
		return errors.New("no suitable locations found")
	}
	if len(l) > 1 {
		return errors.New("multiple locations found, try to specify")
	}

	err = a.Storage.SetUserLocation(ctx, uID, l[0])
	if err != nil {
		return err
	}

	return nil
}

func (a *App) GetCurrentUVIndex(ctx context.Context) (float32, error) {
	uID, err := a.getUserIDFromCtx(ctx)
	if err != nil {
		return 0, err
	}
	u, err := a.Storage.GetUserSettingsOrCreate(ctx, uID)
	if err != nil {
		return 0, err
	}

	uv, err := a.WeatherAPI.GetCurrentUVIndex(ctx, u.Location)
	if err != nil {
		return 0, err
	}
	return uv, nil
}

func (a *App) getUserIDFromCtx(ctx context.Context) (int64, error) {
	u, ok := ctx.Value(UserIDCtxKey).(int64)
	if !ok {
		return 0, fmt.Errorf("can't convert %s value (%v) to int64", UserIDCtxKey, u)
	}

	return u, nil
}

