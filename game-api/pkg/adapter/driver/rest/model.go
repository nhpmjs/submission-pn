package rest

import (
	"time"

	"github.com/google/uuid"
)

type NewPlayer struct {
	Name string `json:"name"`
}

type PlayerRes struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
}

type NewGame struct {
	Name    string   `json:"name"`
	Players []string `json:"players"`
}

type GetGamePlay struct {
	GameID uuid.UUID `param:"gameId"`
}

type SubmitScoreReq struct {
	GamePlayID uuid.UUID `param:"gameId"`
	PlayerID   uuid.UUID `param:"playerId"`
	Score      int       `json:"score"`
}

type Score struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	Score      int       `json:"score"`
	Frame      int       `json:"frame"`
	Roll       int       `json:"roll"`
	PlayerID   uuid.UUID `json:"playerId"`
	PlayerName string    `json:"playername"`
}

type GamePlay struct {
	ID           uuid.UUID         `json:"id"`
	CreatedAt    time.Time         `json:"createdAt"`
	Status       string            `json:"status"`
	Name         string            `json:"name"`
	CurrentFrame int               `json:"currentFrame"`
	CurrentRoll  int               `json:"currentRoll"`
	CurrentUser  int               `json:"currentUser"`
	Participants []GameParticipant `json:"participants"`
}

type GameParticipant struct {
	PlayerID uuid.UUID `json:"playerId"`
	Name     string    `json:"name"`
}
