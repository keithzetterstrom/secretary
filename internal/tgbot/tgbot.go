package tgbot

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	intererrors "github.com/keithzetterstrom/secretary/internal/errors"
	"github.com/keithzetterstrom/secretary/internal/pkg/user"
	"github.com/keithzetterstrom/secretary/internal/repository/lru"
	"github.com/keithzetterstrom/secretary/utils/logger"
)

type Config struct {
	Token string `yaml:"token"`
}

type service struct {
	bot      *tgbotapi.BotAPI
	ctx      context.Context
	user     user.Usecase
	cfg      Config
	sessions FSM
	log      logger.Logger
}

type TgBot interface {
	Execute()
}

func New(
	ctx context.Context,
	lru lru.LRU,
	user user.Usecase,
	cfg Config,
	log logger.Logger,
) (TgBot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to start new bot api")
	}
	bot.Debug = false

	_, err = bot.Request(
		tgbotapi.NewSetMyCommands(
			tgbotapi.BotCommand{Command: "/start", Description: "start"},
			tgbotapi.BotCommand{Command: "/activate", Description: "activate"},
			tgbotapi.BotCommand{Command: "/get", Description: "get"},
		),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to do bot api request")
	}

	sessions, err := NewFsm(ctx, lru)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create new fsm")
	}

	return &service{
		bot:      bot,
		ctx:      ctx,
		user:     user,
		cfg:      cfg,
		sessions: sessions,
		log:      log,
	}, nil
}

func (s *service) processState(
	msg tgbotapi.MessageConfig,
	data string,
	sentFromId int64,
	state *StateContext,
) (tgbotapi.MessageConfig, error) {
	msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      false,
	}

	msg.Text = ""

	var err error
	switch state.State {
	case FsmStateStart:
		msg.Text = "Привет"
		msg.ReplyMarkup = s.NewMainMenuKeyboard()
	case FsmStateActivate:
		err := s.user.Activate(s.ctx, state.TgUserNameId)
		if err != nil {
			return msg, err
		}
		msg.Text = "Регистрация подтверждена"
		msg.ReplyMarkup = s.NewMainMenuKeyboard()
	case FsmStateGet:
		u, err := s.user.Get(s.ctx, state.TgUserNameId)
		if err != nil {
			return msg, err
		}
		msg.Text = fmt.Sprintf("email: %s\nname: %s", u.Email, u.Name)
		msg.ReplyMarkup = s.NewMainMenuKeyboard()
	default:
		return msg, errors.Wrapf(intererrors.ErrValidation, "failed to process state %d", state.State)
	}

	return msg, err
}

func (s *service) processUpdate(update tgbotapi.Update, replyId int64, data string) (tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(replyId, "")
	msg.ParseMode = "HTML"

	state, _ := s.sessions.Get(update.SentFrom().ID)
	if state == nil {
		// new session
		state = &StateContext{
			State:        FsmStateStart,
			TgId:         update.SentFrom().ID,
			TgUserNameId: update.SentFrom().UserName,
		}

		s.sessions.Create(update.SentFrom().ID, update.SentFrom().UserName)
	}

	switch data {
	case "/start":
		state.State = FsmStateStart
	case "/activate":
		state.State = FsmStateActivate
	case "/get":
		state.State = FsmStateGet
	}

	return s.processState(msg, data, update.SentFrom().ID, state)
}

func (s *service) Execute() {
	s.log.Info(fmt.Sprintf("Starting bot %s", s.bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.bot.GetUpdatesChan(u)

	for update := range updates {
		var replyId int64
		data := ""

		if update.Message != nil {
			replyId = update.Message.Chat.ID
			data = update.Message.Text
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := s.bot.Request(callback); err != nil {
				s.log.Error(fmt.Sprintf("failed to request callback: %s", err.Error()))
				continue
			}
			replyId = update.CallbackQuery.Message.Chat.ID
			data = update.CallbackQuery.Data
		} else {
			continue
		}

		msg, err := s.processUpdate(update, replyId, data)
		if err != nil {
			s.log.Error(fmt.Sprintf("failed to set next state: %s", err.Error()))
			continue
		}

		if _, err := s.bot.Send(msg); err != nil {
			s.log.Error(fmt.Sprintf("failed to send message (%s) for user %d: %s", msg.Text, update.SentFrom().ID, err.Error()))
			continue
		}
	}
}

func getUserName(user *tgbotapi.User) string {
	username := fmt.Sprintf("user-%d", user.ID)
	if len(user.UserName) > 1 {
		username = user.UserName
	}
	if len(user.FirstName) > 1 {
		username = user.FirstName
	}

	return username
}
