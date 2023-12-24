package app

import (
	"context"
	weatherAPI "github.com/ympyst/uvindex-tgbot/weather"
	"log"
	"math"
	"os"
)

type App struct {
	weatherClient *weatherAPI.APIClient
}

func NewApp() *App {
	wc := weatherAPI.NewAPIClient(weatherAPI.NewConfiguration())

	return &App{
		wc,
	}
}

func (a *App) SetLocation(search string) {

}

func (a *App) GetCurrentUVIndex(ctx context.Context) int32 {
	authCtx := a.getAuthCtx(ctx)
	res, _, err := a.weatherClient.APIsApi.RealtimeWeather(authCtx, "Москва", nil)
	if err != nil {
		log.Println(err)
	}

	m, ok := res.(map[string]interface{})
	if !ok {
		log.Println("error converting RealtimeWeather result to map")
	}
	cur, ok := m["current"].(map[string]interface{})
	if !ok {
		log.Println("error converting RealtimeWeather result to map")
	}
	//log.Println(cur)
	uv, ok := cur["uv"].(float64)
	if !ok {
		log.Println("error converting uv to float64")
	}
	return int32(math.Round(uv))
}

func (a *App) getAuthCtx(ctx context.Context) context.Context {
	return context.WithValue(context.Background(), weatherAPI.ContextAPIKey, weatherAPI.APIKey{
		Key: os.Getenv("WEATHER_API_TOKEN"),
	})
}

