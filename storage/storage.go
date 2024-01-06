package storage

import (
	"context"
	"errors"
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
			Id:          userId,
			UVThreshold: model.DefaultUVThreshold,
			State:       model.Start,
		}
	}
	return s.users[userId], nil
}

func (s *Storage) SetUsername(ctx context.Context, userId int64, username string) error {
	u, ok := s.users[userId]
	if !ok {
		return errors.New(fmt.Sprintf("user %d not found", userId))
	}
	u.Username = username
	s.users[userId] = u
	return nil
}

func (s *Storage) SetIsGroup(ctx context.Context, userId int64, isGroup bool) error {
	u, ok := s.users[userId]
	if !ok {
		return errors.New(fmt.Sprintf("user %d not found", userId))
	}
	u.IsGroup = isGroup
	s.users[userId] = u
	return nil
}

func (s *Storage) SaveState(ctx context.Context, state *model.UserState) error {
	if _, ok := s.users[state.Id]; ok {
		s.users[state.Id] = *state
	} else {
		return fmt.Errorf("user %d not found", state.Id)
	}
	return nil
}
