package models

import (
	"time"
)

// Player represents a football player
type Player struct {
	ID           int       `json:"id" db:"id"`
	TeamID       int       `json:"team_id" db:"team_id"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Position     string    `json:"position" db:"position"`
	JerseyNumber *int      `json:"jersey_number,omitempty" db:"jersey_number"`
	Height       *int      `json:"height,omitempty" db:"height"` // in inches
	Weight       *int      `json:"weight,omitempty" db:"weight"` // in pounds
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// PlayerStats represents football statistics for a player in a specific game
type PlayerStats struct {
	ID       int `json:"id" db:"id"`
	PlayerID int `json:"player_id" db:"player_id"`
	GameID   int `json:"game_id" db:"game_id"`
	// Offensive stats
	PassingAttempts      *int `json:"passing_attempts,omitempty" db:"passing_attempts"`
	PassingCompletions   *int `json:"passing_completions,omitempty" db:"passing_completions"`
	PassingYards         *int `json:"passing_yards,omitempty" db:"passing_yards"`
	PassingTouchdowns    *int `json:"passing_touchdowns,omitempty" db:"passing_touchdowns"`
	PassingInterceptions *int `json:"passing_interceptions,omitempty" db:"passing_interceptions"`
	RushingAttempts      *int `json:"rushing_attempts,omitempty" db:"rushing_attempts"`
	RushingYards         *int `json:"rushing_yards,omitempty" db:"rushing_yards"`
	RushingTouchdowns    *int `json:"rushing_touchdowns,omitempty" db:"rushing_touchdowns"`
	ReceivingTargets     *int `json:"receiving_targets,omitempty" db:"receiving_targets"`
	Receptions           *int `json:"receptions,omitempty" db:"receptions"`
	ReceivingYards       *int `json:"receiving_yards,omitempty" db:"receiving_yards"`
	ReceivingTouchdowns  *int `json:"receiving_touchdowns,omitempty" db:"receiving_touchdowns"`
	Fumbles              *int `json:"fumbles,omitempty" db:"fumbles"`
	FumblesLost          *int `json:"fumbles_lost,omitempty" db:"fumbles_lost"`
	// Defensive stats
	Tackles                *int `json:"tackles,omitempty" db:"tackles"`
	SoloTackles            *int `json:"solo_tackles,omitempty" db:"solo_tackles"`
	AssistedTackles        *int `json:"assisted_tackles,omitempty" db:"assisted_tackles"`
	Sacks                  *int `json:"sacks,omitempty" db:"sacks"`
	DefensiveInterceptions *int `json:"defensive_interceptions,omitempty" db:"defensive_interceptions"`
	PassDeflections        *int `json:"pass_deflections,omitempty" db:"pass_deflections"`
	ForcedFumbles          *int `json:"forced_fumbles,omitempty" db:"forced_fumbles"`
	FumbleRecoveries       *int `json:"fumble_recoveries,omitempty" db:"fumble_recoveries"`
	DefensiveTouchdowns    *int `json:"defensive_touchdowns,omitempty" db:"defensive_touchdowns"`
	// Special teams
	FieldGoalsAttempted  *int      `json:"field_goals_attempted,omitempty" db:"field_goals_attempted"`
	FieldGoalsMade       *int      `json:"field_goals_made,omitempty" db:"field_goals_made"`
	ExtraPointsAttempted *int      `json:"extra_points_attempted,omitempty" db:"extra_points_attempted"`
	ExtraPointsMade      *int      `json:"extra_points_made,omitempty" db:"extra_points_made"`
	Punts                *int      `json:"punts,omitempty" db:"punts"`
	PuntYards            *int      `json:"punt_yards,omitempty" db:"punt_yards"`
	KickReturns          *int      `json:"kick_returns,omitempty" db:"kick_returns"`
	KickReturnYards      *int      `json:"kick_return_yards,omitempty" db:"kick_return_yards"`
	KickReturnTouchdowns *int      `json:"kick_return_touchdowns,omitempty" db:"kick_return_touchdowns"`
	PuntReturns          *int      `json:"punt_returns,omitempty" db:"punt_returns"`
	PuntReturnYards      *int      `json:"punt_return_yards,omitempty" db:"punt_return_yards"`
	PuntReturnTouchdowns *int      `json:"punt_return_touchdowns,omitempty" db:"punt_return_touchdowns"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// Request/Response structs for Players
type CreatePlayerRequest struct {
	TeamID       int    `json:"team_id" validate:"required"`
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	Position     string `json:"position" validate:"required"`
	JerseyNumber *int   `json:"jersey_number,omitempty"`
	Height       *int   `json:"height,omitempty"`
	Weight       *int   `json:"weight,omitempty"`
}

type UpdatePlayerRequest struct {
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	Position     *string `json:"position,omitempty"`
	JerseyNumber *int    `json:"jersey_number,omitempty"`
	Height       *int    `json:"height,omitempty"`
	Weight       *int    `json:"weight,omitempty"`
}
