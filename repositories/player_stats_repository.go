package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"sports-backend/models"
)

// PlayerStatsRepository defines the interface for player stats data operations
type PlayerStatsRepository interface {
	GetByID(id int) (*models.PlayerStats, error)
	GetAll() ([]*models.PlayerStats, error)
	GetByPlayerID(playerID int) ([]*models.PlayerStats, error)
	GetByGameID(gameID int) ([]*models.PlayerStats, error)
	GetByPlayerAndGame(playerID, gameID int) (*models.PlayerStats, error)
	Create(stats *models.PlayerStats) error
	Update(stats *models.PlayerStats) error
	Delete(id int) error
	Exists(id int) (bool, error)
	ExistsByPlayerAndGame(playerID, gameID int) (bool, error)
}

// playerStatsRepository implements PlayerStatsRepository interface
type playerStatsRepository struct {
	db *sql.DB
}

// NewPlayerStatsRepository creates a new player stats repository
func NewPlayerStatsRepository(db *sql.DB) PlayerStatsRepository {
	return &playerStatsRepository{db: db}
}

// GetByID retrieves player stats by ID
func (r *playerStatsRepository) GetByID(id int) (*models.PlayerStats, error) {
	query := `
		SELECT ps.id, ps.player_id, ps.game_id,
		       ps.passing_attempts, ps.passing_completions, ps.passing_yards, ps.passing_touchdowns, ps.passing_interceptions,
		       ps.rushing_attempts, ps.rushing_yards, ps.rushing_touchdowns,
		       ps.receiving_targets, ps.receptions, ps.receiving_yards, ps.receiving_touchdowns,
		       ps.fumbles, ps.fumbles_lost,
		       ps.tackles, ps.solo_tackles, ps.assisted_tackles, ps.sacks, ps.defensive_interceptions,
		       ps.pass_deflections, ps.forced_fumbles, ps.fumble_recoveries, ps.defensive_touchdowns,
		       ps.field_goals_attempted, ps.field_goals_made, ps.extra_points_attempted, ps.extra_points_made,
		       ps.punts, ps.punt_yards, ps.kick_returns, ps.kick_return_yards, ps.kick_return_touchdowns,
		       ps.punt_returns, ps.punt_return_yards, ps.punt_return_touchdowns,
		       ps.created_at, ps.updated_at,
		       p.first_name, p.last_name, p.position, p.jersey_number,
		       t.name as team_name, t.city as team_city
		FROM player_stats ps
		JOIN players p ON ps.player_id = p.id
		JOIN teams t ON p.team_id = t.id
		WHERE ps.id = ?
	`

	var stats models.PlayerStats
	var firstName, lastName, position, teamName, teamCity string
	var jerseyNumber *int

	err := r.db.QueryRow(query, id).Scan(
		&stats.ID, &stats.PlayerID, &stats.GameID,
		&stats.PassingAttempts, &stats.PassingCompletions, &stats.PassingYards, &stats.PassingTouchdowns, &stats.PassingInterceptions,
		&stats.RushingAttempts, &stats.RushingYards, &stats.RushingTouchdowns,
		&stats.ReceivingTargets, &stats.Receptions, &stats.ReceivingYards, &stats.ReceivingTouchdowns,
		&stats.Fumbles, &stats.FumblesLost,
		&stats.Tackles, &stats.SoloTackles, &stats.AssistedTackles, &stats.Sacks, &stats.DefensiveInterceptions,
		&stats.PassDeflections, &stats.ForcedFumbles, &stats.FumbleRecoveries, &stats.DefensiveTouchdowns,
		&stats.FieldGoalsAttempted, &stats.FieldGoalsMade, &stats.ExtraPointsAttempted, &stats.ExtraPointsMade,
		&stats.Punts, &stats.PuntYards, &stats.KickReturns, &stats.KickReturnYards, &stats.KickReturnTouchdowns,
		&stats.PuntReturns, &stats.PuntReturnYards, &stats.PuntReturnTouchdowns,
		&stats.CreatedAt, &stats.UpdatedAt,
		&firstName, &lastName, &position, &jerseyNumber, &teamName, &teamCity,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("player stats with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get player stats: %w", err)
	}

	return &stats, nil
}

// GetAll retrieves all player stats
func (r *playerStatsRepository) GetAll() ([]*models.PlayerStats, error) {
	query := `
		SELECT ps.id, ps.player_id, ps.game_id,
		       ps.passing_attempts, ps.passing_completions, ps.passing_yards, ps.passing_touchdowns, ps.passing_interceptions,
		       ps.rushing_attempts, ps.rushing_yards, ps.rushing_touchdowns,
		       ps.receiving_targets, ps.receptions, ps.receiving_yards, ps.receiving_touchdowns,
		       ps.fumbles, ps.fumbles_lost,
		       ps.tackles, ps.solo_tackles, ps.assisted_tackles, ps.sacks, ps.defensive_interceptions,
		       ps.pass_deflections, ps.forced_fumbles, ps.fumble_recoveries, ps.defensive_touchdowns,
		       ps.field_goals_attempted, ps.field_goals_made, ps.extra_points_attempted, ps.extra_points_made,
		       ps.punts, ps.punt_yards, ps.kick_returns, ps.kick_return_yards, ps.kick_return_touchdowns,
		       ps.punt_returns, ps.punt_return_yards, ps.punt_return_touchdowns,
		       ps.created_at, ps.updated_at,
		       p.first_name, p.last_name, p.position, p.jersey_number,
		       t.name as team_name, t.city as team_city
		FROM player_stats ps
		JOIN players p ON ps.player_id = p.id
		JOIN teams t ON p.team_id = t.id
		ORDER BY ps.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query player stats: %w", err)
	}
	defer rows.Close()

	var statsList []*models.PlayerStats
	for rows.Next() {
		var stats models.PlayerStats
		var firstName, lastName, position, teamName, teamCity string
		var jerseyNumber *int

		err := rows.Scan(
			&stats.ID, &stats.PlayerID, &stats.GameID,
			&stats.PassingAttempts, &stats.PassingCompletions, &stats.PassingYards, &stats.PassingTouchdowns, &stats.PassingInterceptions,
			&stats.RushingAttempts, &stats.RushingYards, &stats.RushingTouchdowns,
			&stats.ReceivingTargets, &stats.Receptions, &stats.ReceivingYards, &stats.ReceivingTouchdowns,
			&stats.Fumbles, &stats.FumblesLost,
			&stats.Tackles, &stats.SoloTackles, &stats.AssistedTackles, &stats.Sacks, &stats.DefensiveInterceptions,
			&stats.PassDeflections, &stats.ForcedFumbles, &stats.FumbleRecoveries, &stats.DefensiveTouchdowns,
			&stats.FieldGoalsAttempted, &stats.FieldGoalsMade, &stats.ExtraPointsAttempted, &stats.ExtraPointsMade,
			&stats.Punts, &stats.PuntYards, &stats.KickReturns, &stats.KickReturnYards, &stats.KickReturnTouchdowns,
			&stats.PuntReturns, &stats.PuntReturnYards, &stats.PuntReturnTouchdowns,
			&stats.CreatedAt, &stats.UpdatedAt,
			&firstName, &lastName, &position, &jerseyNumber, &teamName, &teamCity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan player stats: %w", err)
		}
		statsList = append(statsList, &stats)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating player stats: %w", err)
	}

	return statsList, nil
}

// GetByPlayerID retrieves all stats for a specific player
func (r *playerStatsRepository) GetByPlayerID(playerID int) ([]*models.PlayerStats, error) {
	query := `
		SELECT ps.id, ps.player_id, ps.game_id,
		       ps.passing_attempts, ps.passing_completions, ps.passing_yards, ps.passing_touchdowns, ps.passing_interceptions,
		       ps.rushing_attempts, ps.rushing_yards, ps.rushing_touchdowns,
		       ps.receiving_targets, ps.receptions, ps.receiving_yards, ps.receiving_touchdowns,
		       ps.fumbles, ps.fumbles_lost,
		       ps.tackles, ps.solo_tackles, ps.assisted_tackles, ps.sacks, ps.defensive_interceptions,
		       ps.pass_deflections, ps.forced_fumbles, ps.fumble_recoveries, ps.defensive_touchdowns,
		       ps.field_goals_attempted, ps.field_goals_made, ps.extra_points_attempted, ps.extra_points_made,
		       ps.punts, ps.punt_yards, ps.kick_returns, ps.kick_return_yards, ps.kick_return_touchdowns,
		       ps.punt_returns, ps.punt_return_yards, ps.punt_return_touchdowns,
		       ps.created_at, ps.updated_at,
		       p.first_name, p.last_name, p.position, p.jersey_number,
		       t.name as team_name, t.city as team_city
		FROM player_stats ps
		JOIN players p ON ps.player_id = p.id
		JOIN teams t ON p.team_id = t.id
		WHERE ps.player_id = ?
		ORDER BY ps.created_at DESC
	`

	rows, err := r.db.Query(query, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query player stats by player: %w", err)
	}
	defer rows.Close()

	var statsList []*models.PlayerStats
	for rows.Next() {
		var stats models.PlayerStats
		var firstName, lastName, position, teamName, teamCity string
		var jerseyNumber *int

		err := rows.Scan(
			&stats.ID, &stats.PlayerID, &stats.GameID,
			&stats.PassingAttempts, &stats.PassingCompletions, &stats.PassingYards, &stats.PassingTouchdowns, &stats.PassingInterceptions,
			&stats.RushingAttempts, &stats.RushingYards, &stats.RushingTouchdowns,
			&stats.ReceivingTargets, &stats.Receptions, &stats.ReceivingYards, &stats.ReceivingTouchdowns,
			&stats.Fumbles, &stats.FumblesLost,
			&stats.Tackles, &stats.SoloTackles, &stats.AssistedTackles, &stats.Sacks, &stats.DefensiveInterceptions,
			&stats.PassDeflections, &stats.ForcedFumbles, &stats.FumbleRecoveries, &stats.DefensiveTouchdowns,
			&stats.FieldGoalsAttempted, &stats.FieldGoalsMade, &stats.ExtraPointsAttempted, &stats.ExtraPointsMade,
			&stats.Punts, &stats.PuntYards, &stats.KickReturns, &stats.KickReturnYards, &stats.KickReturnTouchdowns,
			&stats.PuntReturns, &stats.PuntReturnYards, &stats.PuntReturnTouchdowns,
			&stats.CreatedAt, &stats.UpdatedAt,
			&firstName, &lastName, &position, &jerseyNumber, &teamName, &teamCity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan player stats: %w", err)
		}
		statsList = append(statsList, &stats)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating player stats: %w", err)
	}

	return statsList, nil
}

// GetByGameID retrieves all stats for a specific game
func (r *playerStatsRepository) GetByGameID(gameID int) ([]*models.PlayerStats, error) {
	query := `
		SELECT ps.id, ps.player_id, ps.game_id,
		       ps.passing_attempts, ps.passing_completions, ps.passing_yards, ps.passing_touchdowns, ps.passing_interceptions,
		       ps.rushing_attempts, ps.rushing_yards, ps.rushing_touchdowns,
		       ps.receiving_targets, ps.receptions, ps.receiving_yards, ps.receiving_touchdowns,
		       ps.fumbles, ps.fumbles_lost,
		       ps.tackles, ps.solo_tackles, ps.assisted_tackles, ps.sacks, ps.defensive_interceptions,
		       ps.pass_deflections, ps.forced_fumbles, ps.fumble_recoveries, ps.defensive_touchdowns,
		       ps.field_goals_attempted, ps.field_goals_made, ps.extra_points_attempted, ps.extra_points_made,
		       ps.punts, ps.punt_yards, ps.kick_returns, ps.kick_return_yards, ps.kick_return_touchdowns,
		       ps.punt_returns, ps.punt_return_yards, ps.punt_return_touchdowns,
		       ps.created_at, ps.updated_at,
		       p.first_name, p.last_name, p.position, p.jersey_number,
		       t.name as team_name, t.city as team_city
		FROM player_stats ps
		JOIN players p ON ps.player_id = p.id
		JOIN teams t ON p.team_id = t.id
		WHERE ps.game_id = ?
		ORDER BY t.name ASC, p.last_name ASC, p.first_name ASC
	`

	rows, err := r.db.Query(query, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to query player stats by game: %w", err)
	}
	defer rows.Close()

	var statsList []*models.PlayerStats
	for rows.Next() {
		var stats models.PlayerStats
		var firstName, lastName, position, teamName, teamCity string
		var jerseyNumber *int

		err := rows.Scan(
			&stats.ID, &stats.PlayerID, &stats.GameID,
			&stats.PassingAttempts, &stats.PassingCompletions, &stats.PassingYards, &stats.PassingTouchdowns, &stats.PassingInterceptions,
			&stats.RushingAttempts, &stats.RushingYards, &stats.RushingTouchdowns,
			&stats.ReceivingTargets, &stats.Receptions, &stats.ReceivingYards, &stats.ReceivingTouchdowns,
			&stats.Fumbles, &stats.FumblesLost,
			&stats.Tackles, &stats.SoloTackles, &stats.AssistedTackles, &stats.Sacks, &stats.DefensiveInterceptions,
			&stats.PassDeflections, &stats.ForcedFumbles, &stats.FumbleRecoveries, &stats.DefensiveTouchdowns,
			&stats.FieldGoalsAttempted, &stats.FieldGoalsMade, &stats.ExtraPointsAttempted, &stats.ExtraPointsMade,
			&stats.Punts, &stats.PuntYards, &stats.KickReturns, &stats.KickReturnYards, &stats.KickReturnTouchdowns,
			&stats.PuntReturns, &stats.PuntReturnYards, &stats.PuntReturnTouchdowns,
			&stats.CreatedAt, &stats.UpdatedAt,
			&firstName, &lastName, &position, &jerseyNumber, &teamName, &teamCity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan player stats: %w", err)
		}
		statsList = append(statsList, &stats)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating player stats: %w", err)
	}

	return statsList, nil
}

// GetByPlayerAndGame retrieves stats for a specific player in a specific game
func (r *playerStatsRepository) GetByPlayerAndGame(playerID, gameID int) (*models.PlayerStats, error) {
	query := `
		SELECT ps.id, ps.player_id, ps.game_id,
		       ps.passing_attempts, ps.passing_completions, ps.passing_yards, ps.passing_touchdowns, ps.passing_interceptions,
		       ps.rushing_attempts, ps.rushing_yards, ps.rushing_touchdowns,
		       ps.receiving_targets, ps.receptions, ps.receiving_yards, ps.receiving_touchdowns,
		       ps.fumbles, ps.fumbles_lost,
		       ps.tackles, ps.solo_tackles, ps.assisted_tackles, ps.sacks, ps.defensive_interceptions,
		       ps.pass_deflections, ps.forced_fumbles, ps.fumble_recoveries, ps.defensive_touchdowns,
		       ps.field_goals_attempted, ps.field_goals_made, ps.extra_points_attempted, ps.extra_points_made,
		       ps.punts, ps.punt_yards, ps.kick_returns, ps.kick_return_yards, ps.kick_return_touchdowns,
		       ps.punt_returns, ps.punt_return_yards, ps.punt_return_touchdowns,
		       ps.created_at, ps.updated_at,
		       p.first_name, p.last_name, p.position, p.jersey_number,
		       t.name as team_name, t.city as team_city
		FROM player_stats ps
		JOIN players p ON ps.player_id = p.id
		JOIN teams t ON p.team_id = t.id
		WHERE ps.player_id = ? AND ps.game_id = ?
	`

	var stats models.PlayerStats
	var firstName, lastName, position, teamName, teamCity string
	var jerseyNumber *int

	err := r.db.QueryRow(query, playerID, gameID).Scan(
		&stats.ID, &stats.PlayerID, &stats.GameID,
		&stats.PassingAttempts, &stats.PassingCompletions, &stats.PassingYards, &stats.PassingTouchdowns, &stats.PassingInterceptions,
		&stats.RushingAttempts, &stats.RushingYards, &stats.RushingTouchdowns,
		&stats.ReceivingTargets, &stats.Receptions, &stats.ReceivingYards, &stats.ReceivingTouchdowns,
		&stats.Fumbles, &stats.FumblesLost,
		&stats.Tackles, &stats.SoloTackles, &stats.AssistedTackles, &stats.Sacks, &stats.DefensiveInterceptions,
		&stats.PassDeflections, &stats.ForcedFumbles, &stats.FumbleRecoveries, &stats.DefensiveTouchdowns,
		&stats.FieldGoalsAttempted, &stats.FieldGoalsMade, &stats.ExtraPointsAttempted, &stats.ExtraPointsMade,
		&stats.Punts, &stats.PuntYards, &stats.KickReturns, &stats.KickReturnYards, &stats.KickReturnTouchdowns,
		&stats.PuntReturns, &stats.PuntReturnYards, &stats.PuntReturnTouchdowns,
		&stats.CreatedAt, &stats.UpdatedAt,
		&firstName, &lastName, &position, &jerseyNumber, &teamName, &teamCity,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("player stats for player %d in game %d not found", playerID, gameID)
		}
		return nil, fmt.Errorf("failed to get player stats: %w", err)
	}

	return &stats, nil
}

// Create adds new player stats to the database
func (r *playerStatsRepository) Create(stats *models.PlayerStats) error {
	query := `
		INSERT INTO player_stats (
			player_id, game_id,
			passing_attempts, passing_completions, passing_yards, passing_touchdowns, passing_interceptions,
			rushing_attempts, rushing_yards, rushing_touchdowns,
			receiving_targets, receptions, receiving_yards, receiving_touchdowns,
			fumbles, fumbles_lost,
			tackles, solo_tackles, assisted_tackles, sacks, defensive_interceptions,
			pass_deflections, forced_fumbles, fumble_recoveries, defensive_touchdowns,
			field_goals_attempted, field_goals_made, extra_points_attempted, extra_points_made,
			punts, punt_yards, kick_returns, kick_return_yards, kick_return_touchdowns,
			punt_returns, punt_return_yards, punt_return_touchdowns,
			created_at, updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	currentTime := time.Now()
	result, err := r.db.Exec(query,
		stats.PlayerID, stats.GameID,
		stats.PassingAttempts, stats.PassingCompletions, stats.PassingYards, stats.PassingTouchdowns, stats.PassingInterceptions,
		stats.RushingAttempts, stats.RushingYards, stats.RushingTouchdowns,
		stats.ReceivingTargets, stats.Receptions, stats.ReceivingYards, stats.ReceivingTouchdowns,
		stats.Fumbles, stats.FumblesLost,
		stats.Tackles, stats.SoloTackles, stats.AssistedTackles, stats.Sacks, stats.DefensiveInterceptions,
		stats.PassDeflections, stats.ForcedFumbles, stats.FumbleRecoveries, stats.DefensiveTouchdowns,
		stats.FieldGoalsAttempted, stats.FieldGoalsMade, stats.ExtraPointsAttempted, stats.ExtraPointsMade,
		stats.Punts, stats.PuntYards, stats.KickReturns, stats.KickReturnYards, stats.KickReturnTouchdowns,
		stats.PuntReturns, stats.PuntReturnYards, stats.PuntReturnTouchdowns,
		currentTime, currentTime,
	)
	if err != nil {
		return fmt.Errorf("failed to create player stats: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get player stats ID: %w", err)
	}

	stats.ID = int(id)
	stats.CreatedAt = currentTime
	stats.UpdatedAt = currentTime

	return nil
}

// Update modifies existing player stats
func (r *playerStatsRepository) Update(stats *models.PlayerStats) error {
	query := `
		UPDATE player_stats SET
			passing_attempts = ?, passing_completions = ?, passing_yards = ?, passing_touchdowns = ?, passing_interceptions = ?,
			rushing_attempts = ?, rushing_yards = ?, rushing_touchdowns = ?,
			receiving_targets = ?, receptions = ?, receiving_yards = ?, receiving_touchdowns = ?,
			fumbles = ?, fumbles_lost = ?,
			tackles = ?, solo_tackles = ?, assisted_tackles = ?, sacks = ?, defensive_interceptions = ?,
			pass_deflections = ?, forced_fumbles = ?, fumble_recoveries = ?, defensive_touchdowns = ?,
			field_goals_attempted = ?, field_goals_made = ?, extra_points_attempted = ?, extra_points_made = ?,
			punts = ?, punt_yards = ?, kick_returns = ?, kick_return_yards = ?, kick_return_touchdowns = ?,
			punt_returns = ?, punt_return_yards = ?, punt_return_touchdowns = ?,
			updated_at = ?
		WHERE id = ?
	`

	currentTime := time.Now()
	result, err := r.db.Exec(query,
		stats.PassingAttempts, stats.PassingCompletions, stats.PassingYards, stats.PassingTouchdowns, stats.PassingInterceptions,
		stats.RushingAttempts, stats.RushingYards, stats.RushingTouchdowns,
		stats.ReceivingTargets, stats.Receptions, stats.ReceivingYards, stats.ReceivingTouchdowns,
		stats.Fumbles, stats.FumblesLost,
		stats.Tackles, stats.SoloTackles, stats.AssistedTackles, stats.Sacks, stats.DefensiveInterceptions,
		stats.PassDeflections, stats.ForcedFumbles, stats.FumbleRecoveries, stats.DefensiveTouchdowns,
		stats.FieldGoalsAttempted, stats.FieldGoalsMade, stats.ExtraPointsAttempted, stats.ExtraPointsMade,
		stats.Punts, stats.PuntYards, stats.KickReturns, stats.KickReturnYards, stats.KickReturnTouchdowns,
		stats.PuntReturns, stats.PuntReturnYards, stats.PuntReturnTouchdowns,
		currentTime, stats.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update player stats: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("player stats with ID %d not found", stats.ID)
	}

	stats.UpdatedAt = currentTime
	return nil
}

// Delete removes player stats from the database
func (r *playerStatsRepository) Delete(id int) error {
	query := "DELETE FROM player_stats WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete player stats: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("player stats with ID %d not found", id)
	}

	return nil
}

// Exists checks if player stats exist by ID
func (r *playerStatsRepository) Exists(id int) (bool, error) {
	query := "SELECT 1 FROM player_stats WHERE id = ? LIMIT 1"
	var exists int
	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check player stats existence: %w", err)
	}
	return true, nil
}

// ExistsByPlayerAndGame checks if player stats exist for a specific player and game
func (r *playerStatsRepository) ExistsByPlayerAndGame(playerID, gameID int) (bool, error) {
	query := "SELECT 1 FROM player_stats WHERE player_id = ? AND game_id = ? LIMIT 1"
	var exists int
	err := r.db.QueryRow(query, playerID, gameID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check player stats existence: %w", err)
	}
	return true, nil
}
