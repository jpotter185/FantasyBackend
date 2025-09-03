package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"sports-backend/models"
)

// PlayerRepository defines the interface for player data operations
type PlayerRepository interface {
	GetByID(id int) (*models.Player, error)
	GetAll() ([]*models.Player, error)
	GetByTeamID(teamID int) ([]*models.Player, error)
	Create(player *models.Player) error
	Update(player *models.Player) error
	Delete(id int) error
	Exists(id int) (bool, error)
}

// playerRepository implements PlayerRepository interface
type playerRepository struct {
	db *sql.DB
}

// NewPlayerRepository creates a new player repository
func NewPlayerRepository(db *sql.DB) PlayerRepository {
	return &playerRepository{db: db}
}

// GetByID retrieves a player by their ID
func (r *playerRepository) GetByID(id int) (*models.Player, error) {
	query := `
		SELECT p.id, p.team_id, p.first_name, p.last_name, p.position, 
		       p.jersey_number, p.height, p.weight, p.created_at, p.updated_at,
		       t.name as team_name, t.city as team_city
		FROM players p
		JOIN teams t ON p.team_id = t.id
		WHERE p.id = ?
	`

	var player models.Player
	var teamName, teamCity string
	err := r.db.QueryRow(query, id).Scan(
		&player.ID, &player.TeamID, &player.FirstName, &player.LastName, &player.Position,
		&player.JerseyNumber, &player.Height, &player.Weight, &player.CreatedAt, &player.UpdatedAt,
		&teamName, &teamCity,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("player with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get player: %w", err)
	}

	return &player, nil
}

// GetAll retrieves all players
func (r *playerRepository) GetAll() ([]*models.Player, error) {
	query := `
		SELECT p.id, p.team_id, p.first_name, p.last_name, p.position, 
		       p.jersey_number, p.height, p.weight, p.created_at, p.updated_at,
		       t.name as team_name, t.city as team_city
		FROM players p
		JOIN teams t ON p.team_id = t.id
		ORDER BY p.last_name ASC, p.first_name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query players: %w", err)
	}
	defer rows.Close()

	var players []*models.Player
	for rows.Next() {
		var player models.Player
		var teamName, teamCity string
		err := rows.Scan(
			&player.ID, &player.TeamID, &player.FirstName, &player.LastName, &player.Position,
			&player.JerseyNumber, &player.Height, &player.Weight, &player.CreatedAt, &player.UpdatedAt,
			&teamName, &teamCity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan player: %w", err)
		}
		players = append(players, &player)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating players: %w", err)
	}

	return players, nil
}

// GetByTeamID retrieves all players for a specific team
func (r *playerRepository) GetByTeamID(teamID int) ([]*models.Player, error) {
	query := `
		SELECT p.id, p.team_id, p.first_name, p.last_name, p.position, 
		       p.jersey_number, p.height, p.weight, p.created_at, p.updated_at,
		       t.name as team_name, t.city as team_city
		FROM players p
		JOIN teams t ON p.team_id = t.id
		WHERE p.team_id = ?
		ORDER BY p.position ASC, p.jersey_number ASC
	`

	rows, err := r.db.Query(query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to query players by team: %w", err)
	}
	defer rows.Close()

	var players []*models.Player
	for rows.Next() {
		var player models.Player
		var teamName, teamCity string
		err := rows.Scan(
			&player.ID, &player.TeamID, &player.FirstName, &player.LastName, &player.Position,
			&player.JerseyNumber, &player.Height, &player.Weight, &player.CreatedAt, &player.UpdatedAt,
			&teamName, &teamCity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan player: %w", err)
		}
		players = append(players, &player)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating players: %w", err)
	}

	return players, nil
}

// Create adds a new player to the database
func (r *playerRepository) Create(player *models.Player) error {
	query := `
		INSERT INTO players (team_id, first_name, last_name, position, jersey_number, height, weight, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	currentTime := time.Now()
	result, err := r.db.Exec(query,
		player.TeamID, player.FirstName, player.LastName, player.Position,
		player.JerseyNumber, player.Height, player.Weight, currentTime, currentTime,
	)
	if err != nil {
		return fmt.Errorf("failed to create player: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get player ID: %w", err)
	}

	player.ID = int(id)
	player.CreatedAt = currentTime
	player.UpdatedAt = currentTime

	return nil
}

// Update modifies an existing player
func (r *playerRepository) Update(player *models.Player) error {
	query := `
		UPDATE players 
		SET team_id = ?, first_name = ?, last_name = ?, position = ?, 
		    jersey_number = ?, height = ?, weight = ?, updated_at = ?
		WHERE id = ?
	`

	currentTime := time.Now()
	result, err := r.db.Exec(query,
		player.TeamID, player.FirstName, player.LastName, player.Position,
		player.JerseyNumber, player.Height, player.Weight, currentTime, player.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update player: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("player with ID %d not found", player.ID)
	}

	player.UpdatedAt = currentTime
	return nil
}

// Delete removes a player from the database
func (r *playerRepository) Delete(id int) error {
	query := "DELETE FROM players WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete player: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("player with ID %d not found", id)
	}

	return nil
}

// Exists checks if a player exists by ID
func (r *playerRepository) Exists(id int) (bool, error) {
	query := "SELECT 1 FROM players WHERE id = ? LIMIT 1"
	var exists int
	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check player existence: %w", err)
	}
	return true, nil
}
