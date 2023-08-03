package tgbot

import (
	"context"
	"fmt"
	"time"

	"github.com/keithzetterstrom/secretary/internal/repository/lru"
)

type sessionService struct {
	ctx      context.Context
	sessions lru.LRU
}

type Sessions interface {
	Set(session StateContext, userId int64, expiration time.Duration)
	Get(userId int64) *StateContext
	Delete(userId int64)
}

func (s *sessionService) Set(session StateContext, userId int64) {
	s.sessions.Set(fmt.Sprintf("%d", userId), session)
}

func (s *sessionService) Get(userId int64) *StateContext {
	currentState, ok := s.sessions.Get(fmt.Sprintf("%d", userId))
	if !ok {
		return nil
	}

	result := currentState.(StateContext)

	return &result
}

func (s *sessionService) Delete(userId int64) {
	//err := s.sessions.Delete(s.ctx, fmt.Sprintf("%d", userId))
}
