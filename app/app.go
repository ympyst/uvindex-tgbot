package app

import (
	"context"
	"log"
	weatherAPI "github.com/ympyst/uvindex-tgbot/weather"
)

type App struct {
	weatherClient *weatherAPI.APIClient
}

func NewApp() *App {
	wCfg := &weatherAPI.Configuration{
		UserAgent:     "UV Index Telegram bot",
	}
	wc := weatherAPI.NewAPIClient(wCfg)

	return &App{
		wc,
	}
}

func (a *App) SetLocation(search string) {

}

func (a *App) GetCurrentUVIndex(ctx context.Context) int32 {
	res, _, _ := a.weatherClient.APIsApi.RealtimeWeather(ctx, "Москва", nil)
	cur, ok := res.(weatherAPI.Current)
	if (!ok) {
		log.Println("error converting RealtimeWeather result to Current type")
	}
	return cur.Uv
}

