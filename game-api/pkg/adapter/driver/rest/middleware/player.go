package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type PlayerMiddlewareConfig struct {
	Skipper middleware.Skipper
}

func PlayerMiddleware(config PlayerMiddlewareConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper != nil && config.Skipper(c) {
				return next(c)
			}

			playerId := c.Request().Header.Get("player-id")
			if playerId == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Player has not been inited")
			}
			u, err := uuid.Parse(playerId)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid player-id")
			}
			c.Set("playerId", u)
			return next(c)
		}
	}
}
