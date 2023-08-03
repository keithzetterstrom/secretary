package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	handlers "github.com/keithzetterstrom/secretary/internal/http"
	"github.com/keithzetterstrom/secretary/internal/http/registration"
	"github.com/keithzetterstrom/secretary/internal/pkg/user/repository"
	"github.com/keithzetterstrom/secretary/internal/pkg/user/usecase"
	"github.com/keithzetterstrom/secretary/internal/repository/docs"
	"github.com/keithzetterstrom/secretary/internal/repository/lru"
	"github.com/keithzetterstrom/secretary/internal/tgbot"
	"github.com/keithzetterstrom/secretary/internal/validator"
	"github.com/keithzetterstrom/secretary/utils/logger"
	"github.com/keithzetterstrom/secretary/utils/must"
)

type baseAPI struct {
	docsClient  *docs.Client
	validate    *validator.Validator
	echoService *echo.Echo
	cache       lru.LRU
	bot         tgbot.TgBot
	logger      logger.Logger
	cfg         Config
	ctx         context.Context
	cancel      context.CancelFunc
}

func Start() *baseAPI {
	var cfg Config

	err := NewConfig(&cfg)
	must.Must(err)

	l, err := logger.NewLogger(fmt.Sprintf("secretary_%s", cfg.Env))
	must.Must(err)

	defer func() {
		if panicErr := recover(); panicErr != nil {
			l.Error("Panic", zap.Any("error", panicErr))
			os.Exit(1)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())

	validate, err := validator.InitValidate()
	must.Must(err)

	docsClient, err := docs.New(cfg.DocsConfig, l)
	must.Must(err)

	cache := lru.New()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	// e.HTTPErrorHandler = customHTTPErrorHandler

	userRepo := repository.New(docsClient)

	userUsecase := usecase.New(userRepo)

	registrtionHandler := registration.New(validate, userUsecase, l)

	handlers.Router(e, l, registrtionHandler)

	bot, err := tgbot.New(ctx, cache, userUsecase, cfg.BotConfig, l)
	must.Must(err)

	base := baseAPI{
		docsClient:  docsClient,
		validate:    validate,
		echoService: e,
		cache:       cache,
		bot:         bot,
		logger:      l,
		cfg:         cfg,
		ctx:         ctx,
		cancel:      cancel,
	}

	return &base
}

func (b *baseAPI) Close() {
	b.echoService.Close()
	b.cancel()
	b.logger.Info("stop server")
}

func (b *baseAPI) ServerStart() {
	b.logger.Info(
		"start server",
		zap.String("address", fmt.Sprintf("%s:%s", b.cfg.ServiceConfig.Host, b.cfg.ServiceConfig.Port)),
	)

	if err := b.echoService.Start(fmt.Sprintf(
		"%s:%s",
		b.cfg.ServiceConfig.Host,
		b.cfg.ServiceConfig.Port,
	)); err != nil && err != http.ErrServerClosed {
		b.logger.Error("shutting down the server")
	}
}

func (b *baseAPI) BotStart() {
	b.logger.Info("start bot")

	b.bot.Execute()
}

func main() {
	b := Start()
	defer b.Close()

	go b.ServerStart()
	go b.BotStart()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
}
