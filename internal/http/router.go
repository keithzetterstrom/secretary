package http

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/keithzetterstrom/secretary/internal/http/registration"
	"github.com/keithzetterstrom/secretary/utils/logger"
)

func Router(
	e *echo.Echo,
	l logger.Logger,
	registrationHandler *registration.Handler,
) {
	api := e.Group("")

	h := promhttp.Handler()
	e.Any("/metrics", echo.WrapHandler(h))

	api.Use(
		logger.EchoRequestLogger(l),
	)

	api.GET("/reg", registrationHandler.Registration)
}
