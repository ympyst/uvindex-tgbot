package storage

import (
	"fmt"
	"github.com/ympyst/uvindex-tgbot/app"
)

const storageInitialCap = 100

type Storage struct {
	users map[int64]app.UserSettings
}

func NewStorage() *Storage {
	return &Storage{
		users: make(map[int64]app.UserSettings, storageInitialCap),
	}
}

func (s *Storage) GetUserSettings(userId int64) *app.UserSettings {
	u, ok := s.users[userId]
	if !ok {
		return nil
	}
	return &u
}

func (s *Storage) SetUserLocation(userId int64, location app.Location) error {
	u, ok := s.users[userId]
	if !ok {
		return fmt.Errorf("user %d not found", userId)
	}
	u.Location = location
	return nil
}
