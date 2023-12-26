package storage

import (
	"context"
	"fmt"
	"github.com/ympyst/uvindex-tgbot/model"
	"log"
)

const storageInitialCap = 100

type Storage struct {
	users map[int64]model.UserSettings
}

func NewStorage() *Storage {
	return &Storage{
		users: make(map[int64]model.UserSettings, storageInitialCap),
	}
}

func (s *Storage) GetUserSettingsOrCreate(ctx context.Context, userId int64) (model.UserSettings, error) {
	u, ok := s.users[userId]
	if !ok {
		s.users[userId] = model.UserSettings{
			UserID:        userId,
			Location:      model.Location{},
			AlertSchedule: model.AlertSchedule{},
			UVThreshold:   model.DefaultUVThreshold,
		}
	}
	log.Println(s.users)
	return u, nil
}

func (s *Storage) SetUserLocation(ctx context.Context, userId int64, location model.Location) error {
	if u, ok := s.users[userId]; ok {
		u.Location = location
		s.users[userId] = u
	} else {
		return fmt.Errorf("user %d not found", userId)
	}
	log.Printf("%+v\n", s.users)
	return nil
}
