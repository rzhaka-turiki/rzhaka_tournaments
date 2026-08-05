package repository

import (
	"context"

	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/google/uuid"
)

type TeamRequestRepository interface {
	Create(ctx context.Context, req *model.TeamRequest) error
	GetByID(ctx context.Context, ID uuid.UUID) (*model.TeamRequest, error)
	GetByTeam(ctx context.Context, teamID uuid.UUID) ([]model.TeamRequest, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]model.TeamRequest, error)
	Exists(ctx context.Context, teamID, userID uuid.UUID) (bool, error)
	Delete(ctx context.Context, ID uuid.UUID) error
}

type teamRequestRepository struct {
	db DBTX
}

func NewTeamRequestRepository(db DBTX) TeamRequestRepository {
	return &teamRequestRepository{
		db: db,
	}
}

func (r *teamRequestRepository) Create(ctx context.Context, req *model.TeamRequest) error {
	query := `
	INSERT INTO team_requests (
		team_id,
		user_id,
		created_by,
		type,
		expires_at,
	)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING
		id,
		created_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		req.TeamID,
		req.UserID,
		req.CreatedBy,
		req.Type,
		req.ExpiresAt,
	).Scan(
		&req.ID,
		&req.CreatedAt,
	)
}

func (r *teamRequestRepository) Exists(ctx context.Context, teamID, userID uuid.UUID) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM team_requests
		WHERE team_id = $1
		  AND user_id = $2
	)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, teamID, userID).Scan(&exists)
	return exists, err
}

func (r *teamRequestRepository) GetByID(ctx context.Context, ID uuid.UUID) (*model.TeamRequest, error) {
	query := `
	SELECT
		id,
		team_id,
		user_id,
		created_by,
		type,
		expires_at,
		created_at
	FROM team_requests
	WHERE id = $1
	`
	var req model.TeamRequest
	err := r.db.QueryRow(ctx, query, ID).Scan(
		&req.ID,
		&req.TeamID,
		&req.UserID,
		&req.CreatedBy,
		&req.Type,
		&req.ExpiresAt,
		&req.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *teamRequestRepository) Delete(ctx context.Context, ID uuid.UUID) error {
	query := `
	DELETE
	FROM team_requests
	WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query, ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *teamRequestRepository) GetByTeam(ctx context.Context, teamID uuid.UUID) ([]model.TeamRequest, error) {
	query := `
	SELECT
		id,
		team_id,
		user_id,
		created_by,
		type,
		expires_at,
		created_at
	FROM team_requests
	WHERE id = $1
	ORDER BY created_at
	`
	rows, err := r.db.Query(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reqs []model.TeamRequest
	for rows.Next() {
		var req model.TeamRequest
		err := rows.Scan(
			&req.ID,
			&req.TeamID,
			&req.UserID,
			&req.CreatedBy,
			&req.Type,
			&req.ExpiresAt,
			&req.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}
	return reqs, nil
}

func (r *teamRequestRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]model.TeamRequest, error) {
	query := `
	SELECT
		id,
		team_id,
		user_id,
		created_by,
		type,
		expires_at,
		created_at
	FROM team_requests
	WHERE user_id = $1
	ORDER BY created_at
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reqs []model.TeamRequest
	for rows.Next() {
		var req model.TeamRequest
		err := rows.Scan(
			&req.ID,
			&req.TeamID,
			&req.UserID,
			&req.CreatedBy,
			&req.Type,
			&req.ExpiresAt,
			&req.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}
	return reqs, nil
}
