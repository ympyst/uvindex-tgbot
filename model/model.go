package model

const DefaultUVThreshold = 3.0

type UserState struct {
	UserID        int64 `json:"user_id"`
	Location      `json:"location,omitempty"`
	AlertSchedule `json:"alert_schedule,omitempty"`
	UVThreshold   float32 `json:"uv_threshold,omitempty"`
	State         `json:"state"`
}

type State int

const (
	Start State = iota
	WaitingForLocation
	Ready
)

type Location struct {
	Name           string  `json:"name,omitempty"`
	Region         string  `json:"region,omitempty"`
	Country        string  `json:"country,omitempty"`
	Lat            float32 `json:"lat,omitempty"`
	Lon            float32 `json:"lon,omitempty"`
	TzId           string  `json:"tz_id,omitempty"`
	LocaltimeEpoch int32   `json:"localtime_epoch,omitempty"`
	Localtime      string  `json:"localtime,omitempty"`
}

type AlertSchedule struct {
}
