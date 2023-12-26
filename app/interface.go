package app

type Storage interface {
	GetUserSettings(userId int64) *UserSettings
	SetUserLocation(userId int64, location Location) error
}
