package models

import (
	"time"
)

// Team represents a football team
type Team struct {
	ID         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	City       string    `json:"city" db:"city"`
	Conference string    `json:"conference" db:"conference"`
	Division   string    `json:"division" db:"division"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// Game represents a football game/match
type Game struct {
	ID         int       `json:"id" db:"id"`
	HomeTeamID int       `json:"home_team_id" db:"home_team_id"`
	AwayTeamID int       `json:"away_team_id" db:"away_team_id"`
	Season     string    `json:"season" db:"season"`
	Week       int       `json:"week" db:"week"`
	GameDate   time.Time `json:"game_date" db:"game_date"`
	Status     string    `json:"status" db:"status"` // scheduled, in_progress, completed, cancelled
	HomeScore  *int      `json:"home_score,omitempty" db:"home_score"`
	AwayScore  *int      `json:"away_score,omitempty" db:"away_score"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// Request/Response structs for Teams
type CreateTeamRequest struct {
	Name       string `json:"name" validate:"required"`
	City       string `json:"city" validate:"required"`
	Conference string `json:"conference" validate:"required"`
	Division   string `json:"division" validate:"required"`
}

type UpdateTeamRequest struct {
	Name       *string `json:"name,omitempty"`
	City       *string `json:"city,omitempty"`
	Conference *string `json:"conference,omitempty"`
	Division   *string `json:"division,omitempty"`
}
