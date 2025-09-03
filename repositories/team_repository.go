package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"sports-backend/models"
)

// TeamRepository defines the interface for team data operations
type TeamRepository interface {
	GetByID(id int) (*models.Team, error)
	GetAll() ([]*models.Team, error)
	GetByConference(conference string) ([]*models.Team, error)
	GetByDivision(division string) ([]*models.Team, error)
	Create(team *models.Team) error
	Update(team *models.Team) error
	Delete(id int) error
	Exists(id int) (bool, error)
}

// teamRepository implements TeamRepository interface
type teamRepository struct {
	db *sql.DB
}

// NewTeamRepository creates a new team repository
func NewTeamRepository(db *sql.DB) TeamRepository {
	return &teamRepository{db: db}
}

// GetByID retrieves a team by their ID
func (r *teamRepository) GetByID(id int) (*models.Team, error) {
	query := `
		SELECT id, name, city, conference, division, created_at, updated_at
		FROM teams WHERE id = ?
	`

	var team models.Team
	err := r.db.QueryRow(query, id).Scan(
		&team.ID, &team.Name, &team.City, &team.Conference,
		&team.Division, &team.CreatedAt, &team.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("team with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return &team, nil
}

// GetAll retrieves all teams
func (r *teamRepository) GetAll() ([]*models.Team, error) {
	query := `
		SELECT id, name, city, conference, division, created_at, updated_at
		FROM teams
		ORDER BY conference ASC, division ASC, name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query teams: %w", err)
	}
	defer rows.Close()

	var teams []*models.Team
	for rows.Next() {
		var team models.Team
		err := rows.Scan(
			&team.ID, &team.Name, &team.City, &team.Conference,
			&team.Division, &team.CreatedAt, &team.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan team: %w", err)
		}
		teams = append(teams, &team)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating teams: %w", err)
	}

	return teams, nil
}

// GetByConference retrieves all teams in a specific conference
func (r *teamRepository) GetByConference(conference string) ([]*models.Team, error) {
	query := `
		SELECT id, name, city, conference, division, created_at, updated_at
		FROM teams
		WHERE conference = ?
		ORDER BY division ASC, name ASC
	`

	rows, err := r.db.Query(query, conference)
	if err != nil {
		return nil, fmt.Errorf("failed to query teams by conference: %w", err)
	}
	defer rows.Close()

	var teams []*models.Team
	for rows.Next() {
		var team models.Team
		err := rows.Scan(
			&team.ID, &team.Name, &team.City, &team.Conference,
			&team.Division, &team.CreatedAt, &team.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan team: %w", err)
		}
		teams = append(teams, &team)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating teams: %w", err)
	}

	return teams, nil
}

// GetByDivision retrieves all teams in a specific division
func (r *teamRepository) GetByDivision(division string) ([]*models.Team, error) {
	query := `
		SELECT id, name, city, conference, division, created_at, updated_at
		FROM teams
		WHERE division = ?
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query, division)
	if err != nil {
		return nil, fmt.Errorf("failed to query teams by division: %w", err)
	}
	defer rows.Close()

	var teams []*models.Team
	for rows.Next() {
		var team models.Team
		err := rows.Scan(
			&team.ID, &team.Name, &team.City, &team.Conference,
			&team.Division, &team.CreatedAt, &team.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan team: %w", err)
		}
		teams = append(teams, &team)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating teams: %w", err)
	}

	return teams, nil
}

// Create adds a new team to the database
func (r *teamRepository) Create(team *models.Team) error {
	query := `
		INSERT INTO teams (name, city, conference, division, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	currentTime := time.Now()
	result, err := r.db.Exec(query,
		team.Name, team.City, team.Conference, team.Division, currentTime, currentTime,
	)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get team ID: %w", err)
	}

	team.ID = int(id)
	team.CreatedAt = currentTime
	team.UpdatedAt = currentTime

	return nil
}

// Update modifies an existing team
func (r *teamRepository) Update(team *models.Team) error {
	query := `
		UPDATE teams 
		SET name = ?, city = ?, conference = ?, division = ?, updated_at = ?
		WHERE id = ?
	`

	currentTime := time.Now()
	result, err := r.db.Exec(query,
		team.Name, team.City, team.Conference, team.Division, currentTime, team.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update team: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("team with ID %d not found", team.ID)
	}

	team.UpdatedAt = currentTime
	return nil
}

// Delete removes a team from the database
func (r *teamRepository) Delete(id int) error {
	query := "DELETE FROM teams WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("team with ID %d not found", id)
	}

	return nil
}

// Exists checks if a team exists by ID
func (r *teamRepository) Exists(id int) (bool, error) {
	query := "SELECT 1 FROM teams WHERE id = ? LIMIT 1"
	var exists int
	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check team existence: %w", err)
	}
	return true, nil
}
