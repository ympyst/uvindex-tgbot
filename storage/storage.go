package storage

import (
	"context"
	"fmt"
	"github.com/ympyst/uvindex-tgbot/model"
)

const storageInitialCap = 100

type Storage struct {
	users map[int64]model.UserState
}

func NewStorage() *Storage {
	return &Storage{
		users: make(map[int64]model.UserState, storageInitialCap),
	}
}

func (s *Storage) GetUserSettingsOrCreate(ctx context.Context, userId int64) (model.UserState, error) {
	_, ok := s.users[userId]
	if !ok {
		s.users[userId] = model.UserState{
			UserID:        userId,
			UVThreshold:   model.DefaultUVThreshold,
			State:         model.Start,
		}
	}
	return s.users[userId], nil
}

func (s *Storage) SaveState(ctx context.Context, state *model.UserState) error {
	if _, ok := s.users[state.UserID]; ok {
		s.users[state.UserID] = *state
	} else {
		return fmt.Errorf("user %d not found", state.UserID)
	}
	return nil
}
