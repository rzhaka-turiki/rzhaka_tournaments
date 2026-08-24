package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
)

const (
	getInviteLinkByIDQuery = `
		SELECT
			id,
			team_id,
			token,
			created_by,
			max_uses,
			used_count,
			expires_at,
			created_at
		FROM team_invite_links
		WHERE id = $1;
	`

	getInviteLinkByTokenQuery = `
		SELECT
			id,
			team_id,
			token,
			created_by,
			max_uses,
			used_count,
			expires_at,
			created_at
		FROM team_invite_links
		WHERE token = $1;
	`

	incrementInviteLinkUsageQuery = `
		UPDATE team_invite_links
		SET used_count = used_count + 1
		WHERE id = $1
			AND used_count < max_uses
		RETURNING used_count;
	`
	deleteInviteLinkQuery = `
		DELETE
		FROM team_invite_links
		WHERE id = $1;
	`

	createInviteLinkQuery = `
	INSERT INTO team_invite_links (
    	team_id,
    	token,
    	created_by,
    	max_uses,
    	expires_at
	)
	VALUES (
    	$1,$2,$3,$4,$5
	)
		RETURNING id, created_at, used_count;
	`

	getByTeamQuery = `
	SELECT
		id,
		team_id,
		token,
		created_by,
		max_uses,
		used_count,
		expires_at,
		created_at
	FROM team_invite_links
	WHERE team_id = $1
	ORDER BY created_at DESC;
	`

	deleteExpiredByTeamQuery = `
	DELETE
	FROM team_invite_links
	WHERE team_id = $1
		AND expires_at <= NOW();
	`
)

type TeamInviteLinkRepository interface {
	Create(ctx context.Context, req *model.TeamInviteLink) error
	GetByToken(ctx context.Context, token string) (*model.TeamInviteLink, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.TeamInviteLink, error)
	IncrementUsage(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteExpiredByTeam(ctx context.Context, teamID uuid.UUID) error
	GetByTeam(ctx context.Context, teamID uuid.UUID) ([]model.TeamInviteLink, error)
}

type teamInviteLinkRepository struct {
	db DBTX
}

func NewTeamInviteLinkRepository(db DBTX) TeamInviteLinkRepository {
	return &teamInviteLinkRepository{
		db: db,
	}
}

func (r *teamInviteLinkRepository) Create(ctx context.Context, req *model.TeamInviteLink) error {
	return r.db.QueryRow(ctx, createInviteLinkQuery,
		req.TeamID,
		req.Token,
		req.CreatedBy,
		req.MaxUses,
		req.ExpiresAt,
	).Scan(&req.ID, &req.CreatedAt, &req.UsedCount)
}

func (r *teamInviteLinkRepository) GetByToken(ctx context.Context, token string) (*model.TeamInviteLink, error) {
	row := r.db.QueryRow(ctx, getInviteLinkByTokenQuery, token)
	return scanInviteLink(row)
}

func (r *teamInviteLinkRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.TeamInviteLink, error) {
	row := r.db.QueryRow(ctx, getInviteLinkByIDQuery, id)
	return scanInviteLink(row)
}

func (r *teamInviteLinkRepository) IncrementUsage(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, incrementInviteLinkUsageQuery, id)
	return err
}

func (r *teamInviteLinkRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, deleteInviteLinkQuery, id)
	return err
}

func scanInviteLink(row pgx.Row) (*model.TeamInviteLink, error) {
	var link model.TeamInviteLink

	err := row.Scan(
		&link.ID,
		&link.TeamID,
		&link.Token,
		&link.CreatedBy,
		&link.MaxUses,
		&link.UsedCount,
		&link.ExpiresAt,
		&link.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (r *teamInviteLinkRepository) GetByTeam(ctx context.Context, teamID uuid.UUID) ([]model.TeamInviteLink, error) {
	rows, err := r.db.Query(ctx, getByTeamQuery, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	links := make([]model.TeamInviteLink, 0)
	for rows.Next() {
		link, err := scanInviteLink(rows)
		if err != nil {
			return nil, err
		}
		links = append(links, *link)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return links, nil
}

func (r *teamInviteLinkRepository) DeleteExpiredByTeam(ctx context.Context, teamID uuid.UUID) error {
	_, err := r.db.Exec(ctx, deleteExpiredByTeamQuery, teamID)
	if err != nil {
		return err
	}
	return nil
}
