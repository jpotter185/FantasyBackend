package repositories

import (
	"database/sql"
	"fmt"
	"sports-backend/models"
	"time"
)

// GameRepository defines the interface for game data operations
type GameRepository interface {
	GetAll() ([]*models.Game, error)
	GetByID(id int) (*models.Game, error)
	Create(game *models.Game) error
	Update(game *models.Game) error
	Delete(id int) error
	GetByTeamID(teamID int) ([]*models.Game, error)
	GetBySeason(season string) ([]*models.Game, error)
	GetByWeek(season string, week int) ([]*models.Game, error)
	Exists(id int) (bool, error)
}

// gameRepository implements the GameRepository interface
type gameRepository struct {
	db *sql.DB
}

// NewGameRepository creates a new game repository
func NewGameRepository(db *sql.DB) GameRepository {
	return &gameRepository{db: db}
}

// GetAll retrieves all games with team information
func (r *gameRepository) GetAll() ([]*models.Game, error) {
	query := `
		SELECT 
			g.id, g.home_team_id, g.away_team_id, g.season, g.week, 
			g.game_date, g.status, g.home_score, g.away_score, 
			g.created_at, g.updated_at,
			ht.name as home_team_name, ht.city as home_team_city,
			at.name as away_team_name, at.city as away_team_city
		FROM games g
		JOIN teams ht ON g.home_team_id = ht.id
		JOIN teams at ON g.away_team_id = at.id
		ORDER BY g.game_date DESC, g.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query games: %w", err)
	}
	defer rows.Close()

	var games []*models.Game
	for rows.Next() {
		var game models.Game
		var homeTeamName, homeTeamCity, awayTeamName, awayTeamCity string

		err := rows.Scan(
			&game.ID, &game.HomeTeamID, &game.AwayTeamID, &game.Season, &game.Week,
			&game.GameDate, &game.Status, &game.HomeScore, &game.AwayScore,
			&game.CreatedAt, &game.UpdatedAt,
			&homeTeamName, &homeTeamCity, &awayTeamName, &awayTeamCity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan game: %w", err)
		}

		games = append(games, &game)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating games: %w", err)
	}

	return games, nil
}

// GetByID retrieves a game by ID with team information
func (r *gameRepository) GetByID(id int) (*models.Game, error) {
	query := `
		SELECT 
			g.id, g.home_team_id, g.away_team_id, g.season, g.week, 
			g.game_date, g.status, g.home_score, g.away_score, 
			g.created_at, g.updated_at,
			ht.name as home_team_name, ht.city as home_team_city,
			at.name as away_team_name, at.city as away_team_city
		FROM games g
		JOIN teams ht ON g.home_team_id = ht.id
		JOIN teams at ON g.away_team_id = at.id
		WHERE g.id = ?
	`

	var game models.Game
	var homeTeamName, homeTeamCity, awayTeamName, awayTeamCity string

	err := r.db.QueryRow(query, id).Scan(
		&game.ID, &game.HomeTeamID, &game.AwayTeamID, &game.Season, &game.Week,
		&game.GameDate, &game.Status, &game.HomeScore, &game.AwayScore,
		&game.CreatedAt, &game.UpdatedAt,
		&homeTeamName, &homeTeamCity, &awayTeamName, &awayTeamCity,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("game with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	return &game, nil
}

// Create creates a new game
func (r *gameRepository) Create(game *models.Game) error {
	query := `
		INSERT INTO games (
			home_team_id, away_team_id, season, week, game_date, status, 
			home_score, away_score, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	currentTime := time.Now()
	result, err := r.db.Exec(query,
		game.HomeTeamID, game.AwayTeamID, game.Season, game.Week,
		game.GameDate, game.Status, game.HomeScore, game.AwayScore,
		currentTime, currentTime,
	)

	if err != nil {
		return fmt.Errorf("failed to create game: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get game ID: %w", err)
	}

	game.ID = int(id)
	game.CreatedAt = currentTime
	game.UpdatedAt = currentTime

	return nil
}

// Update updates an existing game
func (r *gameRepository) Update(game *models.Game) error {
	query := `
		UPDATE games SET 
			home_team_id = ?, away_team_id = ?, season = ?, week = ?, 
			game_date = ?, status = ?, home_score = ?, away_score = ?, 
			updated_at = ?
		WHERE id = ?
	`

	currentTime := time.Now()
	result, err := r.db.Exec(query,
		game.HomeTeamID, game.AwayTeamID, game.Season, game.Week,
		game.GameDate, game.Status, game.HomeScore, game.AwayScore,
		currentTime, game.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("game with ID %d not found", game.ID)
	}

	game.UpdatedAt = currentTime
	return nil
}

// Delete deletes a game by ID
func (r *gameRepository) Delete(id int) error {
	query := `DELETE FROM games WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete game: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("game with ID %d not found", id)
	}

	return nil
}

// GetByTeamID retrieves all games for a specific team (both home and away)
func (r *gameRepository) GetByTeamID(teamID int) ([]*models.Game, error) {
	query := `
		SELECT 
			g.id, g.home_team_id, g.away_team_id, g.season, g.week, 
			g.game_date, g.status, g.home_score, g.away_score, 
			g.created_at, g.updated_at,
			ht.name as home_team_name, ht.city as home_team_city,
			at.name as away_team_name, at.city as away_team_city
		FROM games g
		JOIN teams ht ON g.home_team_id = ht.id
		JOIN teams at ON g.away_team_id = at.id
		WHERE g.home_team_id = ? OR g.away_team_id = ?
		ORDER BY g.game_date DESC, g.created_at DESC
	`

	rows, err := r.db.Query(query, teamID, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to query games by team: %w", err)
	}
	defer rows.Close()

	var games []*models.Game
	for rows.Next() {
		var game models.Game
		var homeTeamName, homeTeamCity, awayTeamName, awayTeamCity string

		err := rows.Scan(
			&game.ID, &game.HomeTeamID, &game.AwayTeamID, &game.Season, &game.Week,
			&game.GameDate, &game.Status, &game.HomeScore, &game.AwayScore,
			&game.CreatedAt, &game.UpdatedAt,
			&homeTeamName, &homeTeamCity, &awayTeamName, &awayTeamCity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan game: %w", err)
		}

		games = append(games, &game)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating games: %w", err)
	}

	return games, nil
}

// GetBySeason retrieves all games for a specific season
func (r *gameRepository) GetBySeason(season string) ([]*models.Game, error) {
	query := `
		SELECT 
			g.id, g.home_team_id, g.away_team_id, g.season, g.week, 
			g.game_date, g.status, g.home_score, g.away_score, 
			g.created_at, g.updated_at,
			ht.name as home_team_name, ht.city as home_team_city,
			at.name as away_team_name, at.city as away_team_city
		FROM games g
		JOIN teams ht ON g.home_team_id = ht.id
		JOIN teams at ON g.away_team_id = at.id
		WHERE g.season = ?
		ORDER BY g.week ASC, g.game_date ASC
	`

	rows, err := r.db.Query(query, season)
	if err != nil {
		return nil, fmt.Errorf("failed to query games by season: %w", err)
	}
	defer rows.Close()

	var games []*models.Game
	for rows.Next() {
		var game models.Game
		var homeTeamName, homeTeamCity, awayTeamName, awayTeamCity string

		err := rows.Scan(
			&game.ID, &game.HomeTeamID, &game.AwayTeamID, &game.Season, &game.Week,
			&game.GameDate, &game.Status, &game.HomeScore, &game.AwayScore,
			&game.CreatedAt, &game.UpdatedAt,
			&homeTeamName, &homeTeamCity, &awayTeamName, &awayTeamCity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan game: %w", err)
		}

		games = append(games, &game)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating games: %w", err)
	}

	return games, nil
}

// GetByWeek retrieves all games for a specific week in a season
func (r *gameRepository) GetByWeek(season string, week int) ([]*models.Game, error) {
	query := `
		SELECT 
			g.id, g.home_team_id, g.away_team_id, g.season, g.week, 
			g.game_date, g.status, g.home_score, g.away_score, 
			g.created_at, g.updated_at,
			ht.name as home_team_name, ht.city as home_team_city,
			at.name as away_team_name, at.city as away_team_city
		FROM games g
		JOIN teams ht ON g.home_team_id = ht.id
		JOIN teams at ON g.away_team_id = at.id
		WHERE g.season = ? AND g.week = ?
		ORDER BY g.game_date ASC
	`

	rows, err := r.db.Query(query, season, week)
	if err != nil {
		return nil, fmt.Errorf("failed to query games by week: %w", err)
	}
	defer rows.Close()

	var games []*models.Game
	for rows.Next() {
		var game models.Game
		var homeTeamName, homeTeamCity, awayTeamName, awayTeamCity string

		err := rows.Scan(
			&game.ID, &game.HomeTeamID, &game.AwayTeamID, &game.Season, &game.Week,
			&game.GameDate, &game.Status, &game.HomeScore, &game.AwayScore,
			&game.CreatedAt, &game.UpdatedAt,
			&homeTeamName, &homeTeamCity, &awayTeamName, &awayTeamCity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan game: %w", err)
		}

		games = append(games, &game)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating games: %w", err)
	}

	return games, nil
}

// Exists checks if a game exists by ID
func (r *gameRepository) Exists(id int) (bool, error) {
	query := `SELECT 1 FROM games WHERE id = ? LIMIT 1`

	var exists int
	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if game exists: %w", err)
	}

	return true, nil
}
