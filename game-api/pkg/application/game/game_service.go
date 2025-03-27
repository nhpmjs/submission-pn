package game

import (
	"context"
	"megastrongbow-api/pkg/application/dto"
	"megastrongbow-api/pkg/application/port"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type gameService struct {
	datastore port.GameDatastore
	validator *validator.Validate
}

// NewGame implements port.GameService.
func (g *gameService) NewGame(ctx context.Context, gameplay dto.NewGamePlay) (*dto.GamePlay, error) {
	return g.datastore.NewGame(ctx, gameplay)
}

// GetScore implements port.GameService.
func (g *gameService) GetScore(ctx context.Context, id uuid.UUID) (*dto.GaemScore, error) {
	return g.datastore.GetScore(ctx, id)
}

func (g *gameService) StartGame(ctx context.Context, score dto.SubmitScoreReq2) error {
	return g.datastore.StartGame(ctx, score)
}

// SubmitScore implements port.GameService.
func (g *gameService) SubmitScore(ctx context.Context, score dto.SubmitScoreReq) (*dto.GamePlay, error) {
	if err := g.validator.Struct(score); err != nil {
		return nil, err
	}

	return g.datastore.SubmitScore(ctx, score)
}

// GetPlayer implements port.GameService.
func (g *gameService) GetPlayer(ctx context.Context, id uuid.UUID) (*dto.Player, error) {
	return g.datastore.GetPlayer(ctx, id)
}

func NewGameService(
	datastore port.GameDatastore,
	validator *validator.Validate,
) *gameService {
	return &gameService{
		datastore: datastore,
		validator: validator,
	}
}

func (g *gameService) NewPlayer(ctx context.Context, name string) (*dto.Player, error) {
	val := struct {
		Name string `validate:"min=1"`
	}{
		Name: name,
	}

	if err := g.validator.Struct(val); err != nil {
		return nil, err
	}

	return g.datastore.NewPlayer(ctx, name)
}

func (g *gameService) AddPlayer(ctx context.Context, player *dto.AddPlayerReq) (*dto.GamePlay, error) {
	return g.datastore.AddPlayer(ctx, player)
}

func (g *gameService) CreateGamePlay(ctx context.Context, name string, owner uuid.UUID) (*dto.GamePlay, error) {
	val := struct {
		Name string `validate:"min=1"`
	}{
		Name: name,
	}

	if err := g.validator.Struct(val); err != nil {
		return nil, err
	}

	return g.datastore.CreateGamePlay(ctx, name, owner)
}

func (d *gameService) GetGamePlay(ctx context.Context, id uuid.UUID) (*dto.GamePlay, error) {
	return d.datastore.GetGamePlay(ctx, id)
}

var _ port.GameService = (*gameService)(nil)
