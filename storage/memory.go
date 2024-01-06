package storage

import (
	"context"
	"fmt"
	"github.com/ympyst/uvindex-tgbot/model"
)

const storageInitialCap = 100

type Memory struct {
	users map[int64]model.UserState
}

func NewMemory() *Memory {
	return &Memory{
		users: make(map[int64]model.UserState, storageInitialCap),
	}
}

func (s *Memory) GetUserStateOrCreate(ctx context.Context, userId int64) (model.UserState, error) {
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

func (s *Memory) SaveState(ctx context.Context, state *model.UserState) error {
	if _, ok := s.users[state.Id]; ok {
		s.users[state.Id] = *state
	} else {
		return fmt.Errorf("user %d not found", state.Id)
	}
	return nil
}
