package model

const DefaultUVThreshold = 3.0

type UserState struct {
	Id            int64  `bson:"id"`
	IsGroup       bool   `bson:"is_group"`
	Username      string `bson:"username"`
	Location      `bson:"location,omitempty"`
	AlertSchedule `bson:"alert_schedule,omitempty"`
	UVThreshold   float32 `bson:"uv_threshold,omitempty"`
	State         `bson:"state"`
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
