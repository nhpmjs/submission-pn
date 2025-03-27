package dto

import (
	"time"

	"github.com/google/uuid"
)

type GamePlay struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	Status       string
	Name         string
	CurrentFrame int
	CurrentRoll  int
	CurrentUser  int
	Participants []GameParticipant
}

type GameParticipant struct {
	PlayerID uuid.UUID
	Name     string
	JoinedAt time.Time
}

type AddPlayerReq struct {
	GamePlayID uuid.UUID
	PlayerID   uuid.UUID
}

type Player struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Name      string
}

type SubmitScoreReq struct {
	GamePlayID uuid.UUID
	PlayerID   uuid.UUID
	Score      int
	Frame      int
	Roll       int
}

type SubmitScoreReq2 struct {
	GamePlayID uuid.UUID
	PlayerID   uuid.UUID
	Score      string `validate:"required,score"`
	Frame      int    `validate:"required,min=1,max=10"`
}

type GaemScore struct {
	Scores []Score
}

type Score struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Score     int
	Frame     int
	Roll      int
	Player    *Player
}

type NewGamePlay struct {
	Name    string   `json:"name"`
	Players []string `json:"players"`
}
