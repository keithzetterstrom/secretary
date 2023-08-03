package tgbot

import "github.com/keithzetterstrom/secretary/internal/pkg/models"

// State machine states

const (
	FsmStateStart    = 1
	FsmStateActivate = 2
	FsmStateSuccess  = 3
	FsmStateGet      = 4
)

// Contexts

type StateContext struct {
	State        uint8
	TgId         int64
	TgUserNameId string
	User         models.UserRegistration
}
