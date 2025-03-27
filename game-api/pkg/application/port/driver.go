package port

import (
	"context"
	"megastrongbow-api/pkg/application/dto"

	"github.com/google/uuid"
)

type GameService interface {
	NewGame(ctx context.Context, gameplay dto.NewGamePlay) (*dto.GamePlay, error)
	CreateGamePlay(ctx context.Context, name string, owner uuid.UUID) (*dto.GamePlay, error)
	AddPlayer(ctx context.Context, player *dto.AddPlayerReq) (*dto.GamePlay, error)
	NewPlayer(ctx context.Context, name string) (*dto.Player, error)
	GetPlayer(ctx context.Context, id uuid.UUID) (*dto.Player, error)
	GetGamePlay(ctx context.Context, id uuid.UUID) (*dto.GamePlay, error)

	SubmitScore(ctx context.Context, score dto.SubmitScoreReq) (*dto.GamePlay, error)

	StartGame(ctx context.Context, score dto.SubmitScoreReq2) error
	GetScore(ctx context.Context, id uuid.UUID) (*dto.GaemScore, error)
}
