package tgbot

import (
	"context"

	intererrors "github.com/keithzetterstrom/secretary/internal/errors"
	"github.com/keithzetterstrom/secretary/internal/repository/lru"
)

type fsmService struct {
	sessions sessionService
}

type FSM interface {
	NextState(tgId int64, state uint8) error
	Create(tgId int64, tgUserName string) error
	Get(tgId int64) (*StateContext, error)
	Delete(tgId int64) error
}

func NewFsm(ctx context.Context, redisHandler lru.LRU) (FSM, error) {
	return &fsmService{
		sessions: sessionService{ctx: ctx, sessions: redisHandler},
	}, nil
}

func (s *fsmService) NextState(userId int64, state uint8) error {
	currentState := s.sessions.Get(userId)

	currentState.State = state
	s.sessions.Set(*currentState, userId)

	return nil
}

func (s *fsmService) Create(tgId int64, tgUserName string) error {
	currentState := s.sessions.Get(tgId)
	if currentState != nil {
		currentState.TgId = tgId
		currentState.TgUserNameId = tgUserName
		s.sessions.Set(*currentState, tgId)
		return nil
	}

	session := StateContext{
		TgId:         tgId,
		TgUserNameId: tgUserName,
	}

	s.sessions.Set(session, tgId)

	return nil
}

func (s *fsmService) Get(tgId int64) (*StateContext, error) {
	currentState := s.sessions.Get(tgId)
	if currentState == nil {
		return nil, intererrors.ErrUserNotFound
	}

	return currentState, nil
}

func (s *fsmService) Delete(tgId int64) error {
	/*err := s.sessions.Delete(tgId)
	if err != nil {
		return err
	}*/
	return nil
}
