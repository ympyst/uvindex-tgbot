package app

import (
	"context"
	"fmt"
	"github.com/ympyst/uvindex-tgbot/storage"
	weatherAPI "github.com/ympyst/uvindex-tgbot/weather"
	"log"
	"math"
	"os"
)

const UserIDCtxKey = "UserID"

type App struct {
	weatherClient *weatherAPI.APIClient //TODO interface
	Storage
}

func NewApp() *App {
	wc := weatherAPI.NewAPIClient(weatherAPI.NewConfiguration())

	return &App{
		wc,
		storage.NewStorage(),
	}
}

func (a *App) SetLocation(ctx context.Context, searchQuery string) error {
	_, err := a.getUserIDFromCtx(ctx)
	if err != nil {
		return err
	}

	//authCtx := a.getAuthCtx(ctx)
	//res, _, err := a.weatherClient.APIsApi.SearchAutocompleteWeather(authCtx, searchQuery)
	//if err != nil {
	//	log.Println(err)
	//}

	return nil
}

func (a *App) GetCurrentUVIndex(ctx context.Context) int32 {
	authCtx := a.getAuthCtx(ctx)
	res, _, err := a.weatherClient.APIsApi.RealtimeWeather(authCtx, "Лос Анджелес", nil)
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

func (a *App) getUserIDFromCtx(ctx context.Context) (int64, error) {
	u, ok := ctx.Value(UserIDCtxKey).(int64)
	if !ok {
		return 0, fmt.Errorf("can't convert %s value (%v) to int64", UserIDCtxKey, u)
	}

	return u, nil
}

