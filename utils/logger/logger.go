package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	contextlib "github.com/keithzetterstrom/secretary/utils/context"
	"github.com/keithzetterstrom/secretary/utils/models"
)

type logger struct {
	logger *zap.Logger
}

type shortError struct {
	e error
}

func (pe shortError) Error() string {
	return pe.e.Error()
}

func ShortError(err error) zap.Field {
	return zap.Error(shortError{err})
}

type Logger interface {
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	WithContext(ctx context.Context) Logger
}

// LogsConfig default format from config file
type LogsConfig struct {
	LogToStderr     bool   `yaml:"log_stderr"`
	LogToSyslog     bool   `yaml:"log_syslog"`
	SyslogTransport string `yaml:"syslog_transport"`
	SyslogHost      string `yaml:"syslog_host"`
	SyslogPort      int    `yaml:"syslog_port"`
}

func NewLogger(serviceName string) (Logger, error) {
	zLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	zLogger = zLogger.With(zap.String("service", serviceName))

	return &logger{logger: zLogger}, nil
}

func (l logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l logger) With(fields ...zap.Field) Logger {
	l.logger = l.logger.With(fields...)
	return l
}

func (l logger) WithContext(ctx context.Context) Logger {
	reqID := contextlib.GetRequestID(ctx)
	if reqID != nil {
		l.logger = l.logger.With(zap.String("request_id", *reqID))
	}

	abc := contextlib.GetXAbc(ctx)
	if abc != nil {
		l.logger = l.logger.With(zap.String("abc", *abc))
	}

	return l
}

func EchoRequestLogger(log Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if err := recover(); err != nil {
					log.Error("Panic", zap.Any("error", err))
				}
			}()

			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			fields := []zapcore.Field{
				zap.String("remote_ip", c.RealIP()),
				zap.String("latency", time.Since(start).String()),
				zap.String("host", req.Host),
				zap.String("request", fmt.Sprintf("%s %s", req.Method, req.RequestURI)),
				zap.Int("status", res.Status),
				zap.Int64("size", res.Size),
				zap.String("user_agent", req.UserAgent()),
			}

			id := req.Header.Get(models.RequestIDHeader)
			if id == "" {
				id = res.Header().Get(models.RequestIDHeader)
			}
			if id != "" {
				fields = append(fields, zap.String("request_id", id))
			}

			abc := res.Header().Get(models.AbcHeader)
			if abc == "" {
				abc = req.Header.Get(models.AbcHeader)
			}
			if abc != "" {
				fields = append(fields, zap.String("abc", abc))
			}

			n := res.Status
			switch {
			case n >= 500:
				log.With(zap.Error(err)).Error("Server error", fields...)
			case n >= 400:
				log.With(ShortError(err)).Warn("Client error", fields...)
			case n >= 300:
				log.Info("Redirection", fields...)
			default:
				log.Info("Success", fields...)
			}

			return nil
		}
	}
}
