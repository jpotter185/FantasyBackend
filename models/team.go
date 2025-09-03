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

// Request/Response structs for Games
type CreateGameRequest struct {
	HomeTeamID int       `json:"home_team_id" validate:"required"`
	AwayTeamID int       `json:"away_team_id" validate:"required"`
	Season     string    `json:"season" validate:"required"`
	Week       int       `json:"week" validate:"required,min=1,max=22"`
	GameDate   time.Time `json:"game_date" validate:"required"`
	Status     string    `json:"status,omitempty" validate:"omitempty,oneof=scheduled in_progress completed cancelled"`
	HomeScore  *int      `json:"home_score,omitempty" validate:"omitempty,min=0"`
	AwayScore  *int      `json:"away_score,omitempty" validate:"omitempty,min=0"`
}

type UpdateGameRequest struct {
	HomeTeamID *int       `json:"home_team_id,omitempty"`
	AwayTeamID *int       `json:"away_team_id,omitempty"`
	Season     *string    `json:"season,omitempty"`
	Week       *int       `json:"week,omitempty" validate:"omitempty,min=1,max=22"`
	GameDate   *time.Time `json:"game_date,omitempty"`
	Status     *string    `json:"status,omitempty" validate:"omitempty,oneof=scheduled in_progress completed cancelled"`
	HomeScore  *int       `json:"home_score,omitempty" validate:"omitempty,min=0"`
	AwayScore  *int       `json:"away_score,omitempty" validate:"omitempty,min=0"`
}
