package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func NewServer() *echo.Echo {
	server := echo.New()

	server.HideBanner = true

	server.Use(middleware.RequestID())
	server.Use(middleware.Recover())
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger := log.Logger
			logger = logger.With().Str("id", c.Response().Header().Get(echo.HeaderXRequestID)).Logger()
			ctx := logger.WithContext(c.Request().Context())
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	})
	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool { return c.Request().URL.Path == "/health" },
	}))

	return server
}
