package repository

import (
	"context"

	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/google/uuid"
)

type TeamMemberRepository interface {
	GetMembers(ctx context.Context, teamID uuid.UUID) ([]model.TeamMember, error)
	Create(ctx context.Context, member *model.TeamMember) error
	Remove(ctx context.Context, teamID, userID uuid.UUID) error
	Exists(ctx context.Context, teamID, userID uuid.UUID) (bool, error)
}

type teamMemberRepository struct {
	db DBTX
}

func NewTeamMemberRepository(db DBTX) TeamMemberRepository {
	return &teamMemberRepository{
		db: db,
	}
}

func (r *teamMemberRepository) GetMembers(ctx context.Context, teamID uuid.UUID) ([]model.TeamMember, error) {
	query := `
	SELECT
		team_id,
		user_id,
		role,
		joined_at
	FROM team_members
	WHERE team_id = $1
	ORDER BY joined_at ASC
	`
	rows, err := r.db.Query(ctx, query, teamID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var members []model.TeamMember

	for rows.Next() {
		var member model.TeamMember
		err := rows.Scan(&member.TeamID, &member.UserID, &member.Role, &member.JoinedAt)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}

func (r *teamMemberRepository) Create(ctx context.Context, member *model.TeamMember) error {
	query := `
	INSERT INTO team_members (
		team_id,
		user_id,
		role
	)
	VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(
		ctx,
		query,
		member.TeamID,
		member.UserID,
		member.Role,
	)

	return err
}

func (r *teamMemberRepository) Exists(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM team_members
		WHERE team_id = $1
		  AND user_id = $2
	)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, teamID, userID).Scan(&exists)
	return exists, err
}

func (r *teamMemberRepository) Remove(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	query := `
	DELETE FROM team_members
	WHERE team_id = $1
	  AND user_id = $2
	`

	result, err := r.db.Exec(ctx, query, teamID, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
