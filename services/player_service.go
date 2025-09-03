package services

import (
	"fmt"
	"strings"

	"sports-backend/models"
	"sports-backend/repositories"
)

// PlayerService defines the interface for player business logic
type PlayerService interface {
	GetPlayer(id int) (*models.Player, error)
	GetAllPlayers() ([]*models.Player, error)
	GetPlayersByTeam(teamID int) ([]*models.Player, error)
	CreatePlayer(req *models.CreatePlayerRequest) (*models.Player, error)
	UpdatePlayer(id int, req *models.UpdatePlayerRequest) (*models.Player, error)
	DeletePlayer(id int) error
}

// playerService implements PlayerService interface
type playerService struct {
	playerRepo repositories.PlayerRepository
	teamRepo   repositories.TeamRepository
}

// NewPlayerService creates a new player service
func NewPlayerService(playerRepo repositories.PlayerRepository, teamRepo repositories.TeamRepository) PlayerService {
	return &playerService{
		playerRepo: playerRepo,
		teamRepo:   teamRepo,
	}
}

// GetPlayer retrieves a player by ID
func (s *playerService) GetPlayer(id int) (*models.Player, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid player ID: %d", id)
	}

	player, err := s.playerRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get player: %w", err)
	}

	return player, nil
}

// GetAllPlayers retrieves all players
func (s *playerService) GetAllPlayers() ([]*models.Player, error) {
	players, err := s.playerRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get players: %w", err)
	}

	return players, nil
}

// GetPlayersByTeam retrieves all players for a specific team
func (s *playerService) GetPlayersByTeam(teamID int) ([]*models.Player, error) {
	if teamID <= 0 {
		return nil, fmt.Errorf("invalid team ID: %d", teamID)
	}

	// Verify team exists
	exists, err := s.teamRepo.Exists(teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify team existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("team with ID %d not found", teamID)
	}

	players, err := s.playerRepo.GetByTeamID(teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get players by team: %w", err)
	}

	return players, nil
}

// CreatePlayer creates a new player
func (s *playerService) CreatePlayer(req *models.CreatePlayerRequest) (*models.Player, error) {
	// Validate request
	if err := s.validateCreatePlayerRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Verify team exists
	exists, err := s.teamRepo.Exists(req.TeamID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify team existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("team with ID %d not found", req.TeamID)
	}

	// Check if jersey number is already taken by another player on the same team
	if req.JerseyNumber != nil {
		players, err := s.playerRepo.GetByTeamID(req.TeamID)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing players: %w", err)
		}

		for _, player := range players {
			if player.JerseyNumber != nil && *player.JerseyNumber == *req.JerseyNumber {
				return nil, fmt.Errorf("jersey number %d is already taken by another player on this team", *req.JerseyNumber)
			}
		}
	}

	// Create player
	player := &models.Player{
		TeamID:       req.TeamID,
		FirstName:    strings.TrimSpace(req.FirstName),
		LastName:     strings.TrimSpace(req.LastName),
		Position:     strings.TrimSpace(req.Position),
		JerseyNumber: req.JerseyNumber,
		Height:       req.Height,
		Weight:       req.Weight,
	}

	if err := s.playerRepo.Create(player); err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	return player, nil
}

// UpdatePlayer updates an existing player
func (s *playerService) UpdatePlayer(id int, req *models.UpdatePlayerRequest) (*models.Player, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid player ID: %d", id)
	}

	// Validate request
	if err := s.validateUpdatePlayerRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing player
	player, err := s.playerRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get player: %w", err)
	}

	// Update fields if provided
	if req.FirstName != nil {
		player.FirstName = strings.TrimSpace(*req.FirstName)
	}
	if req.LastName != nil {
		player.LastName = strings.TrimSpace(*req.LastName)
	}
	if req.Position != nil {
		player.Position = strings.TrimSpace(*req.Position)
	}
	if req.JerseyNumber != nil {
		// Check if jersey number is already taken by another player on the same team
		players, err := s.playerRepo.GetByTeamID(player.TeamID)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing players: %w", err)
		}

		for _, existingPlayer := range players {
			if existingPlayer.ID != id && existingPlayer.JerseyNumber != nil && *existingPlayer.JerseyNumber == *req.JerseyNumber {
				return nil, fmt.Errorf("jersey number %d is already taken by another player on this team", *req.JerseyNumber)
			}
		}
		player.JerseyNumber = req.JerseyNumber
	}
	if req.Height != nil {
		player.Height = req.Height
	}
	if req.Weight != nil {
		player.Weight = req.Weight
	}

	// Update player
	if err := s.playerRepo.Update(player); err != nil {
		return nil, fmt.Errorf("failed to update player: %w", err)
	}

	return player, nil
}

// DeletePlayer deletes a player
func (s *playerService) DeletePlayer(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid player ID: %d", id)
	}

	// Check if player exists
	exists, err := s.playerRepo.Exists(id)
	if err != nil {
		return fmt.Errorf("failed to check player existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("player with ID %d not found", id)
	}

	// TODO: Add business logic here if needed
	// For example: check if player has stats, prevent deletion if they do

	if err := s.playerRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete player: %w", err)
	}

	return nil
}

// validateCreatePlayerRequest validates the create player request
func (s *playerService) validateCreatePlayerRequest(req *models.CreatePlayerRequest) error {
	if req.TeamID <= 0 {
		return fmt.Errorf("team ID is required and must be positive")
	}

	if strings.TrimSpace(req.FirstName) == "" {
		return fmt.Errorf("first name is required")
	}

	if strings.TrimSpace(req.LastName) == "" {
		return fmt.Errorf("last name is required")
	}

	if strings.TrimSpace(req.Position) == "" {
		return fmt.Errorf("position is required")
	}

	// Validate jersey number if provided
	if req.JerseyNumber != nil {
		if *req.JerseyNumber < 0 || *req.JerseyNumber > 99 {
			return fmt.Errorf("jersey number must be between 0 and 99")
		}
	}

	// Validate height if provided
	if req.Height != nil {
		if *req.Height < 60 || *req.Height > 90 { // 5'0" to 7'6"
			return fmt.Errorf("height must be between 60 and 90 inches")
		}
	}

	// Validate weight if provided
	if req.Weight != nil {
		if *req.Weight < 150 || *req.Weight > 400 { // 150 to 400 pounds
			return fmt.Errorf("weight must be between 150 and 400 pounds")
		}
	}

	return nil
}

// validateUpdatePlayerRequest validates the update player request
func (s *playerService) validateUpdatePlayerRequest(req *models.UpdatePlayerRequest) error {
	// Check if at least one field is being updated
	if req.FirstName == nil && req.LastName == nil && req.Position == nil &&
		req.JerseyNumber == nil && req.Height == nil && req.Weight == nil {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Validate individual fields if provided
	if req.FirstName != nil && strings.TrimSpace(*req.FirstName) == "" {
		return fmt.Errorf("first name cannot be empty")
	}

	if req.LastName != nil && strings.TrimSpace(*req.LastName) == "" {
		return fmt.Errorf("last name cannot be empty")
	}

	if req.Position != nil && strings.TrimSpace(*req.Position) == "" {
		return fmt.Errorf("position cannot be empty")
	}

	// Validate jersey number if provided
	if req.JerseyNumber != nil {
		if *req.JerseyNumber < 0 || *req.JerseyNumber > 99 {
			return fmt.Errorf("jersey number must be between 0 and 99")
		}
	}

	// Validate height if provided
	if req.Height != nil {
		if *req.Height < 60 || *req.Height > 90 {
			return fmt.Errorf("height must be between 60 and 90 inches")
		}
	}

	// Validate weight if provided
	if req.Weight != nil {
		if *req.Weight < 150 || *req.Weight > 400 {
			return fmt.Errorf("weight must be between 150 and 400 pounds")
		}
	}

	return nil
}
