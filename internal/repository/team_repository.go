package repository

import (
	"context"

	"github.com/a1uka/rzhaka_tournaments/internal/model"
	"github.com/google/uuid"
)

type TeamRepository interface {
	Create(ctx context.Context, team *model.Team) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Team, error)
	GetOwnedTeams(ctx context.Context, userID uuid.UUID, includeArchived bool) ([]model.Team, error)
	GetMemberTeams(ctx context.Context, userID uuid.UUID, includeArchived bool) ([]model.Team, error)
	CountActiveTeams(ctx context.Context, ownerID uuid.UUID) (int, error)
	IsOwner(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (bool, error)
	Archive(ctx context.Context, teamID uuid.UUID) error
	Restore(ctx context.Context, teamID uuid.UUID) error
	Update(ctx context.Context, team *model.Team) error
}

type teamRepository struct {
	db DBTX
}

func NewTeamRepository(db DBTX) TeamRepository {
	return &teamRepository{
		db: db,
	}
}

func (r *teamRepository) Create(ctx context.Context, team *model.Team) error {
	query := `
	INSERT INTO teams (
		name,
		short_name,
		logo_path,
		logo_dark_path,
		owner_id
	)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING
		id,
		created_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		team.Name,
		team.ShortName,
		team.LogoPath,
		team.LogoDarkPath,
		team.OwnerID,
	).Scan(
		&team.ID,
		&team.CreatedAt,
	)
}

func (r *teamRepository) GetByID(ctx context.Context, ID uuid.UUID) (*model.Team, error) {
	query := `
	SELECT
		id,
		name,
		short_name,
		logo_path,
		logo_dark_path,
		owner_id,
		created_at,
		archived_at
	FROM teams
	WHERE id = $1
	`

	var team model.Team

	err := r.db.QueryRow(ctx, query, ID).Scan(
		&team.ID,
		&team.Name,
		&team.ShortName,
		&team.LogoPath,
		&team.LogoDarkPath,
		&team.OwnerID,
		&team.CreatedAt,
		&team.ArchivedAt,
	)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) AddMember(ctx context.Context, teamID uuid.UUID, member model.TeamMember) error {
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
		teamID,
		member.UserID,
		member.Role,
	)
	return err
}

func (r *teamRepository) GetMembers(ctx context.Context, teamID uuid.UUID) ([]model.TeamMember, error) {
	query := `
	SELECT
		user_id,
		role,
		joined_at
	FROM team_members
	WHERE team_id = $1
	ORDER BY joined_at
	`
	rows, err := r.db.Query(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []model.TeamMember
	for rows.Next() {
		var member model.TeamMember
		err := rows.Scan(
			&member.UserID,
			&member.Role,
			&member.JoinedAt,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, rows.Err()
}

func (r *teamRepository) GetOwnedTeams(ctx context.Context, userID uuid.UUID, includeArchived bool) ([]model.Team, error) {
	query := `
	SELECT
		id,
		name,
		short_name,
		logo_path,
		logo_dark_path,
		owner_id,
		created_at,
		archived_at
	FROM teams
	WHERE owner_id = $1
	`
	if !includeArchived {
		query += ` AND archived_at IS NULL`
	}
	query += ` ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var teams []model.Team
	for rows.Next() {
		var team model.Team
		err := rows.Scan(
			&team.ID,
			&team.Name,
			&team.ShortName,
			&team.LogoPath,
			&team.LogoDarkPath,
			&team.OwnerID,
			&team.CreatedAt,
			&team.ArchivedAt,
		)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, rows.Err()
}

func (r *teamRepository) GetMemberTeams(ctx context.Context, userID uuid.UUID, includeArchived bool) ([]model.Team, error) {
	query := `
	SELECT
    	t.id,
    	t.name,
    	t.short_name,
    	t.logo_path,
    	t.logo_dark_path,
    	t.owner_id,
    	t.created_at,
    	t.archived_at
	FROM teams t
	JOIN team_members tm
    	ON tm.team_id = t.id
	WHERE tm.user_id = $1
	`
	if !includeArchived {
		query += ` AND t.archived_at IS NULL`
	}
	query += ` ORDER BY t.created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var teams []model.Team
	for rows.Next() {
		var team model.Team
		err := rows.Scan(
			&team.ID,
			&team.Name,
			&team.ShortName,
			&team.LogoPath,
			&team.LogoDarkPath,
			&team.OwnerID,
			&team.CreatedAt,
			&team.ArchivedAt,
		)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, rows.Err()
}

func (r *teamRepository) CountActiveTeams(ctx context.Context, ownerID uuid.UUID) (int, error) {
	query := `
	SELECT COUNT(*)
	FROM teams
	WHERE owner_id = $1
	  AND archived_at IS NULL
	`
	var count int
	err := r.db.QueryRow(ctx, query, ownerID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *teamRepository) Archive(ctx context.Context, teamID uuid.UUID) error {
	query := `
	UPDATE teams
	SET archived_at = NOW()
	WHERE id = $1
	  AND archived_at IS NULL
	`
	result, err := r.db.Exec(ctx, query, teamID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *teamRepository) Restore(ctx context.Context, teamID uuid.UUID) error {
	query := `
	UPDATE teams
	SET archived_at = NULL
	WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query, teamID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *teamRepository) IsOwner(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM teams
		WHERE id = $1
		  AND owner_id = $2
	)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, teamID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *teamRepository) RemoveMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	query := `
	DELETE
	FROM team_members
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

func (r *teamRepository) Update(ctx context.Context, team *model.Team) error {
	query := `
	UPDATE teams
	SET
		name = $1,
		short_name = $2,
		logo_path = $3,
		logo_dark_path = $4
	WHERE id = $5
	`

	result, err := r.db.Exec(
		ctx,
		query,
		team.Name,
		team.ShortName,
		team.LogoPath,
		team.LogoDarkPath,
		team.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
