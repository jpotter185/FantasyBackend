package database

import (
	"fmt"
	"log"
)

// RunMigrations creates all necessary database tables
func RunMigrations() error {
	// First, run table creation migrations
	migrations := []struct {
		name string
		sql  string
	}{
		{"teams", createTeamsTable},
		{"games", createGamesTable},
		{"players", createPlayersTable},
		{"player_stats", createPlayerStatsTable},
	}

	for _, migration := range migrations {
		log.Printf("Running migration: %s", migration.name)
		if _, err := DB.Exec(migration.sql); err != nil {
			return fmt.Errorf("failed to run migration %s: %v", migration.name, err)
		}
		log.Printf("Migration %s completed successfully", migration.name)
	}


	log.Println("All database migrations completed successfully")
	return nil
}

// TableExists checks if a table exists in the database
func TableExists(tableName string) (bool, error) {
	query := `
		SELECT name FROM sqlite_master 
		WHERE type='table' AND name=?
	`

	var name string
	err := DB.QueryRow(query, tableName).Scan(&name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if table %s exists: %v", tableName, err)
	}

	return true, nil
}

const createTeamsTable = `
CREATE TABLE IF NOT EXISTS teams (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    city TEXT NOT NULL,
    conference TEXT NOT NULL,
    division TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, city)
);`

const createGamesTable = `
CREATE TABLE IF NOT EXISTS games (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    home_team_id INTEGER NOT NULL,
    away_team_id INTEGER NOT NULL,
    season TEXT NOT NULL,
    week INTEGER NOT NULL,
    game_date DATETIME NOT NULL,
    status TEXT NOT NULL DEFAULT 'scheduled', -- scheduled, in_progress, completed, cancelled
    home_score INTEGER,
    away_score INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (home_team_id) REFERENCES teams (id),
    FOREIGN KEY (away_team_id) REFERENCES teams (id),
    UNIQUE(home_team_id, away_team_id, season, week, game_date)
);`

const createPlayersTable = `
CREATE TABLE IF NOT EXISTS players (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    team_id INTEGER NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    position TEXT NOT NULL,
    jersey_number INTEGER,
    height INTEGER, -- in inches
    weight INTEGER, -- in pounds
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams (id),
    UNIQUE(team_id, first_name, last_name, position, jersey_number)
);`

const createPlayerStatsTable = `
CREATE TABLE IF NOT EXISTS player_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id INTEGER NOT NULL,
    game_id INTEGER NOT NULL,
    
    -- Offensive stats
    passing_attempts INTEGER DEFAULT 0,
    passing_completions INTEGER DEFAULT 0,
    passing_yards INTEGER DEFAULT 0,
    passing_touchdowns INTEGER DEFAULT 0,
    passing_interceptions INTEGER DEFAULT 0,
    
    rushing_attempts INTEGER DEFAULT 0,
    rushing_yards INTEGER DEFAULT 0,
    rushing_touchdowns INTEGER DEFAULT 0,
    
    receiving_targets INTEGER DEFAULT 0,
    receptions INTEGER DEFAULT 0,
    receiving_yards INTEGER DEFAULT 0,
    receiving_touchdowns INTEGER DEFAULT 0,
    
    fumbles INTEGER DEFAULT 0,
    fumbles_lost INTEGER DEFAULT 0,
    
    -- Defensive stats
    tackles INTEGER DEFAULT 0,
    solo_tackles INTEGER DEFAULT 0,
    assisted_tackles INTEGER DEFAULT 0,
    sacks INTEGER DEFAULT 0,
    defensive_interceptions INTEGER DEFAULT 0,
    pass_deflections INTEGER DEFAULT 0,
    forced_fumbles INTEGER DEFAULT 0,
    fumble_recoveries INTEGER DEFAULT 0,
    defensive_touchdowns INTEGER DEFAULT 0,
    
    -- Special teams
    field_goals_attempted INTEGER DEFAULT 0,
    field_goals_made INTEGER DEFAULT 0,
    extra_points_attempted INTEGER DEFAULT 0,
    extra_points_made INTEGER DEFAULT 0,
    punts INTEGER DEFAULT 0,
    punt_yards INTEGER DEFAULT 0,
    kick_returns INTEGER DEFAULT 0,
    kick_return_yards INTEGER DEFAULT 0,
    kick_return_touchdowns INTEGER DEFAULT 0,
    punt_returns INTEGER DEFAULT 0,
    punt_return_yards INTEGER DEFAULT 0,
    punt_return_touchdowns INTEGER DEFAULT 0,
    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (player_id) REFERENCES players (id),
    FOREIGN KEY (game_id) REFERENCES games (id),
    
    -- Ensure one stat record per player per game
    UNIQUE(player_id, game_id)
);`
