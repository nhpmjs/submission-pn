package rest

import (
	"errors"
	"megastrongbow-api/pkg/application/dto"
	"megastrongbow-api/pkg/application/port"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

type gameHandler struct {
	gameService port.GameService
}

func NewGameHandler(
	gameService port.GameService,
) *gameHandler {
	return &gameHandler{
		gameService: gameService,
	}
}

func (h *gameHandler) NewPlayer(c echo.Context) error {
	ctx := c.Request().Context()
	logger := log.Ctx(ctx)
	var req NewPlayer
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	// TODO: prevent multiple submit
	p, err := h.gameService.NewPlayer(ctx, req.Name)
	if err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			logger.Err(err).Msg("badRequest")
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
		}
		logger.Err(err).Msg("internalServerError")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	}

	return c.JSON(200, PlayerRes{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		Name:      p.Name,
	})
}

func (h *gameHandler) NewGame(c echo.Context) error {
	ctx := c.Request().Context()
	logger := log.Ctx(ctx)

	var req NewGame
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	gp, err := h.gameService.NewGame(ctx, dto.NewGamePlay{
		Name:    req.Name,
		Players: req.Players,
	})

	if err != nil {
		logger.Err(err).Msg("internalServerError")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	}

	return c.JSON(200, gp)
}

func (h *gameHandler) GetMe(c echo.Context) error {
	ctx := c.Request().Context()
	logger := log.Ctx(ctx)
	p, err := h.gameService.GetPlayer(ctx, c.Get("playerId").(uuid.UUID))

	if err != nil {
		logger.Err(err).Msg("internalServerError")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	}

	return c.JSON(200, PlayerRes{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		Name:      p.Name,
	})
}

func (h *gameHandler) GetGamePlay(c echo.Context) error {
	ctx := c.Request().Context()
	logger := log.Ctx(ctx)

	var req GetGamePlay
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	p, err := h.gameService.GetGamePlay(ctx, req.GameID)

	if err != nil {
		if errors.Is(err, dto.ErrItemNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Game play not found")
		}
		logger.Err(err).Msg("internalServerError")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	}

	return c.JSON(200, GamePlay{
		ID:           p.ID,
		CreatedAt:    p.CreatedAt,
		Status:       p.Status,
		Name:         p.Name,
		CurrentFrame: p.CurrentFrame,
		CurrentRoll:  p.CurrentRoll,
		CurrentUser:  p.CurrentUser,
		Participants: lo.Map(p.Participants, func(gp dto.GameParticipant, _ int) GameParticipant {
			return GameParticipant{
				PlayerID: gp.PlayerID,
				Name:     gp.Name,
			}
		}),
	})
}

func (h *gameHandler) GetScore(c echo.Context) error {
	ctx := c.Request().Context()
	logger := log.Ctx(ctx)

	var req GetGamePlay
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	p, err := h.gameService.GetScore(ctx, req.GameID)

	if err != nil {
		logger.Err(err).Msg("internalServerError")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	}

	return c.JSON(200, lo.Map(p.Scores, func(s dto.Score, _ int) Score {
		return Score{
			ID:         s.ID,
			CreatedAt:  s.CreatedAt,
			Score:      s.Score,
			Frame:      s.Frame,
			Roll:       s.Roll,
			PlayerID:   s.Player.ID,
			PlayerName: s.Player.Name,
		}
	}))
}

func (h *gameHandler) JoinGame(c echo.Context) error {
	ctx := c.Request().Context()
	logger := log.Ctx(ctx)

	var req GetGamePlay
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	p, err := h.gameService.AddPlayer(ctx, &dto.AddPlayerReq{
		GamePlayID: req.GameID,
		PlayerID:   c.Get("playerId").(uuid.UUID),
	})

	if err != nil {
		if errors.Is(err, dto.ErrDuplicated) {
			return echo.NewHTTPError(http.StatusBadRequest, "Already joined")
		}

		if errors.Is(err, dto.ErrItemNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "Game play or player not found")
		}

		logger.Err(err).Msg("internalServerError")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	}

	return c.JSON(200, PlayerRes{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		Name:      p.Name,
	})
}

func (h *gameHandler) SubmitScore(c echo.Context) error {
	ctx := c.Request().Context()
	logger := log.Ctx(ctx)

	var req SubmitScoreReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	_, err := h.gameService.SubmitScore(ctx, dto.SubmitScoreReq{
		PlayerID:   req.PlayerID,
		GamePlayID: req.GamePlayID,
		Score:      req.Score,
	})

	if err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			logger.Err(err).Msg("badRequest")
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
		}

		if errors.Is(err, dto.ErrDuplicated) {
			return echo.NewHTTPError(http.StatusBadRequest, "Already submitted")
		}

		if errors.Is(err, dto.ErrItemNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "Game not found or has not started")
		}

		logger.Err(err).Msg("internalServerError")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *gameHandler) StartGame(c echo.Context) error {
	ctx := c.Request().Context()
	logger := log.Ctx(ctx)

	var req SubmitScoreReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	// playerId := c.Get("playerId").(uuid.UUID)

	err := h.gameService.StartGame(ctx, dto.SubmitScoreReq2{
		GamePlayID: req.GamePlayID,
	})

	if err != nil {
		if errors.Is(err, dto.ErrGameAlreadyStarted) {
			return echo.NewHTTPError(http.StatusBadRequest, "Already started")
		}
		logger.Err(err).Msg("internalServerError")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	}

	return c.NoContent(http.StatusNoContent)
}
