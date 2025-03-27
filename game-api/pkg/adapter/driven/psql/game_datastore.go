package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"megastrongbow-api/pkg/application/dto"
	"megastrongbow-api/pkg/application/port"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

type gameDatastore struct {
	dbc *bun.DB
}

// NewGame implements port.GameDatastore.
func (d *gameDatastore) NewGame(ctx context.Context, gameplay dto.NewGamePlay) (*dto.GamePlay, error) {
	p := GamePlay{
		Name: gameplay.Name,
	}

	err := d.dbc.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(&p).Exec(ctx); err != nil {
			return fmt.Errorf("unable to create game play: %w", err)
		}

		players := lo.Map(gameplay.Players, func(p string, _ int) *Player {
			return &Player{
				Name: p,
			}
		})

		if _, err := tx.NewInsert().Model(&players).Exec(ctx); err != nil {
			return fmt.Errorf("unable to create players: %w", err)
		}

		participants := lo.Map(players, func(player *Player, _ int) *GameParticipant {
			return &GameParticipant{
				GamePlayID: p.ID,
				PlayerID:   player.ID,
			}
		})

		if _, err := tx.NewInsert().Model(&participants).Exec(ctx); err != nil {
			return fmt.Errorf("unable to add participant %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.GamePlay{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		Status:    p.Status,
		Name:      p.Name,
	}, nil
}

// GetScore implements port.GameDatastore.
func (d *gameDatastore) GetScore(ctx context.Context, id uuid.UUID) (*dto.GaemScore, error) {
	var scores []*Score
	if err := d.dbc.NewSelect().Model(&scores).
		Relation("Player").
		Where("game_play_id = ?", id).
		Order("created_at").
		Scan(ctx); err != nil {
		return nil, err
	}

	return &dto.GaemScore{
		Scores: lo.Map(scores, func(r *Score, _ int) dto.Score {
			return dto.Score{
				ID:        r.ID,
				CreatedAt: r.CreatedAt,
				Score:     r.Score,
				Frame:     r.Frame,
				Roll:      r.Roll,
				Player: &dto.Player{
					ID:        r.PlayerID,
					Name:      r.Player.Name,
					CreatedAt: r.Player.CreatedAt,
				},
			}
		}),
	}, nil
}

func (d *gameDatastore) SubmitScore(ctx context.Context, score dto.SubmitScoreReq) (*dto.GamePlay, error) {
	gp := GamePlay{
		ID: score.GamePlayID,
	}

	if err := d.dbc.NewSelect().Model(&gp).
		Relation("Participants").
		WherePK().
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dto.ErrItemNotFound
		}
		return nil, err
	}

	err := d.dbc.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		playerId := gp.Participants[gp.CurrentUser].PlayerID
		var scores []*Score
		if err := d.dbc.NewSelect().Model(&scores).
			Relation("Player").
			Where("game_play_id = ?", gp.ID).
			Where("player_id = ?", playerId).
			Where("frame = ?", gp.CurrentFrame).
			Scan(ctx); err != nil {
			return fmt.Errorf("unable to get scores: %w", err)
		}

		s := Score{
			GamePlayID: score.GamePlayID,
			PlayerID:   playerId,
			Score:      score.Score,
			Frame:      gp.CurrentFrame,
			Roll:       gp.CurrentRoll,
		}

		if _, err := tx.NewInsert().Model(&s).Exec(ctx); err != nil {
			if err, ok := err.(*pgconn.PgError); ok {
				switch err.ConstraintName {
				case "score_frame_unique":
					return dto.ErrDuplicated
				}
			}
			return err
		}

		scores = append(scores, &s)

		if shouldNextUser(scores, gp.CurrentFrame) {
			gp.nextUser()
		} else {
			gp.nextRoll()
		}

		if gp.CurrentFrame > 10 {
			gp.Status = "done"
		}

		if _, err := tx.NewUpdate().Model(&gp).WherePK().Exec(ctx); err != nil {
			return fmt.Errorf("unable to update game play: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func shouldNextUser(scores []*Score, frame int) bool {
	rolls := len(scores)

	if frame < 10 {
		return rolls == 2 || scores[0].Score == 10
	}

	if rolls == 1 {
		return false
	}

	return rolls == 3 || lo.SumBy(scores, func(s *Score) int { return s.Score }) < 10
}

func increment(n int) int {
	return (n % 2) + 1
}

func (d *gameDatastore) StartGame(ctx context.Context, score dto.SubmitScoreReq2) error {
	p := GamePlay{
		ID: score.GamePlayID,
	}

	if err := d.dbc.NewSelect().Model(&p).
		WherePK().
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.ErrItemNotFound
		}

		return err
	}

	if _, err := d.dbc.NewUpdate().Model(&p).
		WherePK().
		Exec(ctx); err != nil {
		return fmt.Errorf("unable to update game play: %w", err)
	}

	return nil
}

func (d *gameDatastore) GetPlayer(ctx context.Context, id uuid.UUID) (*dto.Player, error) {
	p := Player{
		ID: id,
	}
	if err := d.dbc.NewSelect().Model(&p).WherePK().Scan(ctx); err != nil {
		return nil, err
	}
	return &dto.Player{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		Name:      p.Name,
	}, nil
}

func (d *gameDatastore) NewPlayer(ctx context.Context, name string) (*dto.Player, error) {
	p := Player{
		Name: name,
	}

	if _, err := d.dbc.NewInsert().Model(&p).Exec(ctx); err != nil {
		return nil, err
	}

	return &dto.Player{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		Name:      p.Name,
	}, nil
}

func (d *gameDatastore) AddPlayer(ctx context.Context, player *dto.AddPlayerReq) (*dto.GamePlay, error) {
	p := GameParticipant{
		GamePlayID: player.GamePlayID,
		PlayerID:   player.PlayerID,
	}

	if _, err := d.dbc.NewInsert().Model(&p).Exec(ctx); err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			switch err.ConstraintName {
			case "game_participant_unique":
				return nil, dto.ErrDuplicated

			case "game_participant_player_fk":
			case "game_participant_game_play_fk":
				return nil, dto.ErrItemNotFound
			}
		}
		return nil, err
	}

	return &dto.GamePlay{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
	}, nil
}

func (d *gameDatastore) CreateGamePlay(ctx context.Context, name string, owner uuid.UUID) (*dto.GamePlay, error) {
	p := GamePlay{
		Name: name,
	}

	err := d.dbc.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(&p).Exec(ctx); err != nil {
			return fmt.Errorf("unable to create game play: %w", err)
		}

		p := GameParticipant{
			GamePlayID: p.ID,
			PlayerID:   owner,
		}

		if _, err := tx.NewInsert().Model(&p).Exec(ctx); err != nil {
			return fmt.Errorf("unable to add participant %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.GamePlay{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		Status:    p.Status,
		Name:      p.Name,
	}, nil
}

func (d *gameDatastore) GetGamePlay(ctx context.Context, id uuid.UUID) (*dto.GamePlay, error) {
	p := GamePlay{
		ID: id,
	}

	if err := d.dbc.NewSelect().Model(&p).
		WherePK().
		Relation("Participants.Player").
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dto.ErrItemNotFound
		}

		return nil, err
	}

	return &dto.GamePlay{
		ID:           p.ID,
		CreatedAt:    p.CreatedAt,
		Status:       p.Status,
		Name:         p.Name,
		CurrentFrame: p.CurrentFrame,
		CurrentRoll:  p.CurrentRoll,
		CurrentUser:  p.CurrentUser,
		Participants: lo.Map(p.Participants, func(gp *GameParticipant, _ int) dto.GameParticipant {
			return dto.GameParticipant{
				PlayerID: gp.PlayerID,
				Name:     gp.Player.Name,
				JoinedAt: gp.CreatedAt,
			}
		}),
	}, nil
}

var _ port.GameDatastore = (*gameDatastore)(nil)

func NewStore(dbc *bun.DB) *gameDatastore {
	return &gameDatastore{
		dbc: dbc,
	}
}
