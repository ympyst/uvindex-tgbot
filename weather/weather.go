package weather

import (
	"context"
	"errors"
	"fmt"
	"github.com/ympyst/uvindex-tgbot/model"
	"github.com/ympyst/uvindex-tgbot/weather/swagger"
	"os"
)

const weatherAPITokenEnvKey = "WEATHER_API_TOKEN"

type API struct {
	weatherClient *swagger.APIClient
}

func NewAPI() *API {
	return &API{
		weatherClient: swagger.NewAPIClient(swagger.NewConfiguration()),
	}
}

func (a *API) GetCurrentUVIndex(ctx context.Context, location model.Location) (float32, error)  {
	query := fmt.Sprintf("%s, %s, %s", location.Name, location.Region, location.Country)
	authCtx := a.getAuthCtx(ctx)
	res, _, err := a.weatherClient.APIsApi.RealtimeWeather(authCtx, query, nil)
	if err != nil {
		return 0, err
	}

	m, ok := res.(map[string]interface{})
	if !ok {
		return 0, errors.New("error converting RealtimeWeather result to map")
	}
	cur, ok := m["current"].(map[string]interface{})
	if !ok {
		return 0, errors.New("error converting RealtimeWeather result to map")
	}
	uv, ok := cur["uv"].(float64)
	if !ok {
		return 0, errors.New("error converting uv to float32")
	}
	return float32(uv), nil
}

func (a *API) getAuthCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, swagger.ContextAPIKey, swagger.APIKey{
		Key: os.Getenv(weatherAPITokenEnvKey),
	})
}
