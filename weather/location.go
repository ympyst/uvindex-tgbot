package weather

import (
	"context"
	"errors"
	"github.com/ympyst/uvindex-tgbot/model"
	"github.com/ympyst/uvindex-tgbot/weather/swagger"
	"os"
)

const locationAPITokenEnvKey = "WEATHER_API_TOKEN"

type LocationAPI struct {
	weatherClient *swagger.APIClient
}

func NewLocationAPI() *LocationAPI {
	return &LocationAPI{
		weatherClient: swagger.NewAPIClient(swagger.NewConfiguration()),
	}
}

func (l *LocationAPI) SearchLocationByName(ctx context.Context, searchQuery string) ([]model.Location, error) {
	authCtx := l.getAuthCtx(ctx)
	res, _, err := l.weatherClient.APIsApi.SearchAutocompleteWeather(authCtx, searchQuery)
	if err != nil {
		return nil, err
	}
	resArr, ok := res.([]interface{})
	if !ok {
		return nil, errors.New("can't convert result to array")
	}
	locations := make([]model.Location, len(resArr))
	for i := 0; i < len(resArr); i++ {
		item := resArr[i].(map[string]interface{})
		locations[i] = model.Location{
			Name:           item["name"].(string),
			Region:         item["region"].(string),
			Country:        item["country"].(string),
			Lat:            0,
			Lon:            0,
			TzId:           "",
			LocaltimeEpoch: 0,
			Localtime:      "",
		}
	}
	return locations, nil
}

func (l *LocationAPI) getAuthCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, swagger.ContextAPIKey, swagger.APIKey{
		Key: os.Getenv(locationAPITokenEnvKey),
	})
}