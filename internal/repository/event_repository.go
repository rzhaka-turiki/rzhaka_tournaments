package repository

import (
	"context"

	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

type EventRepository interface {
	Create(ctx context.Context, event *model.Event) error
}

type eventRepository struct {
	db DBTX
}

func NewEventRepository(db DBTX) EventRepository {
	return &eventRepository{
		db: db,
	}
}

func (r *eventRepository) Create(ctx context.Context, event *model.Event) error {
	query := `
		INSERT INTO events (
			actor_id,
			user_id,
			team_id,
			tournament_id,
			registration_id,
			stage_id,
			match_id,
			event_type,
			payload
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9
		)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		event.ActorID,
		event.UserID,
		event.TeamID,
		event.TournamentID,
		event.RegistrationID,
		event.StageID,
		event.MatchID,
		event.EventType,
		event.Payload,
	).Scan(
		&event.ID,
		&event.CreatedAt,
	)
}
