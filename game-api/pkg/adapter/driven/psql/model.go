package psql

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type GameParticipant struct {
	bun.BaseModel `bun:"game_participant,alias:t"`

	ID         uuid.UUID `bun:"id,pk,type:uuid,default:uuid()"`
	CreatedAt  time.Time `bun:"created_at,default:now()"`
	GamePlayID uuid.UUID `bun:"game_play_id,type:uuid"`
	PlayerID   uuid.UUID `bun:"player_id,type:uuid"`

	GamePlay *GamePlay `bun:"join:game_play_id=id,rel:belongs-to"`
	Player   *Player   `bun:"join:player_id=id,rel:belongs-to"`
}

type GamePlay struct {
	bun.BaseModel `bun:"game_play,alias:t"`

	ID           uuid.UUID `bun:"id,pk,type:uuid,default:uuid()"`
	CreatedAt    time.Time `bun:"created_at,default:now()"`
	Status       string    `bun:"status,nullzero"`
	Name         string    `bun:"name"`
	CurrentFrame int       `bun:"current_frame,default:1"`
	CurrentRoll  int       `bun:"current_roll,default:1"`
	CurrentUser  int       `bun:"current_user_index,default:0"`

	Participants []*GameParticipant `bun:"rel:has-many,join:id=game_play_id"`
}

type Player struct {
	bun.BaseModel `bun:"player,alias:t"`

	ID        uuid.UUID `bun:"id,pk,type:uuid,default:uuid()"`
	CreatedAt time.Time `bun:"created_at,default:now()"`
	Name      string    `bun:"name,nullzero"`
}

type Score struct {
	bun.BaseModel `bun:"score,alias:t"`

	ID         uuid.UUID `bun:"id,pk,type:uuid,default:uuid()"`
	CreatedAt  time.Time `bun:"created_at,default:now()"`
	GamePlayID uuid.UUID `bun:"game_play_id,type:uuid,nullzero"`
	PlayerID   uuid.UUID `bun:"player_id,type:uuid,nullzero"`
	Score      int       `bun:"score"`
	Frame      int       `bun:"frame"`
	Roll       int       `bun:"roll"`

	GamePlay *GamePlay `bun:"join:game_play_id=id,rel:belongs-to"`
	Player   *Player   `bun:"join:player_id=id,rel:belongs-to"`
}
