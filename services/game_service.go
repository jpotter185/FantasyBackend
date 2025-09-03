package services

import (
	"fmt"
	"sports-backend/models"
	"sports-backend/repositories"
	"time"
)

// GameService defines the interface for game business logic
type GameService interface {
	GetAllGames() ([]*models.Game, error)
	GetGameByID(id int) (*models.Game, error)
	CreateGame(req *models.CreateGameRequest) (*models.Game, error)
	UpdateGame(id int, req *models.UpdateGameRequest) (*models.Game, error)
	DeleteGame(id int) error
	GetGamesByTeam(teamID int) ([]*models.Game, error)
	GetGamesBySeason(season string) ([]*models.Game, error)
	GetGamesByWeek(season string, week int) ([]*models.Game, error)
}

// gameService implements the GameService interface
type gameService struct {
	gameRepo repositories.GameRepository
	teamRepo repositories.TeamRepository
}

// NewGameService creates a new game service
func NewGameService(gameRepo repositories.GameRepository, teamRepo repositories.TeamRepository) GameService {
	return &gameService{
		gameRepo: gameRepo,
		teamRepo: teamRepo,
	}
}

// GetAllGames retrieves all games
func (s *gameService) GetAllGames() ([]*models.Game, error) {
	return s.gameRepo.GetAll()
}

// GetGameByID retrieves a game by ID
func (s *gameService) GetGameByID(id int) (*models.Game, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid game ID: %d", id)
	}

	return s.gameRepo.GetByID(id)
}

// CreateGame creates a new game
func (s *gameService) CreateGame(req *models.CreateGameRequest) (*models.Game, error) {
	// Validate the request
	if err := s.validateCreateGameRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if both teams exist
	homeTeamExists, err := s.teamRepo.Exists(req.HomeTeamID)
	if err != nil {
		return nil, fmt.Errorf("failed to check home team: %w", err)
	}
	if !homeTeamExists {
		return nil, fmt.Errorf("home team with ID %d not found", req.HomeTeamID)
	}

	awayTeamExists, err := s.teamRepo.Exists(req.AwayTeamID)
	if err != nil {
		return nil, fmt.Errorf("failed to check away team: %w", err)
	}
	if !awayTeamExists {
		return nil, fmt.Errorf("away team with ID %d not found", req.AwayTeamID)
	}

	// Check if teams are different
	if req.HomeTeamID == req.AwayTeamID {
		return nil, fmt.Errorf("home team and away team cannot be the same")
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = "scheduled"
	}

	// Create the game
	game := &models.Game{
		HomeTeamID: req.HomeTeamID,
		AwayTeamID: req.AwayTeamID,
		Season:     req.Season,
		Week:       req.Week,
		GameDate:   req.GameDate,
		Status:     status,
		HomeScore:  req.HomeScore,
		AwayScore:  req.AwayScore,
	}

	if err := s.gameRepo.Create(game); err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	return game, nil
}

// UpdateGame updates an existing game
func (s *gameService) UpdateGame(id int, req *models.UpdateGameRequest) (*models.Game, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid game ID: %d", id)
	}

	// Get the existing game
	game, err := s.gameRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	// Validate the request
	if err := s.validateUpdateGameRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update fields if provided
	if req.HomeTeamID != nil {
		// Check if the new home team exists
		homeTeamExists, err := s.teamRepo.Exists(*req.HomeTeamID)
		if err != nil {
			return nil, fmt.Errorf("failed to check home team: %w", err)
		}
		if !homeTeamExists {
			return nil, fmt.Errorf("home team with ID %d not found", *req.HomeTeamID)
		}
		game.HomeTeamID = *req.HomeTeamID
	}

	if req.AwayTeamID != nil {
		// Check if the new away team exists
		awayTeamExists, err := s.teamRepo.Exists(*req.AwayTeamID)
		if err != nil {
			return nil, fmt.Errorf("failed to check away team: %w", err)
		}
		if !awayTeamExists {
			return nil, fmt.Errorf("away team with ID %d not found", *req.AwayTeamID)
		}
		game.AwayTeamID = *req.AwayTeamID
	}

	// Check if teams are different (after potential updates)
	if game.HomeTeamID == game.AwayTeamID {
		return nil, fmt.Errorf("home team and away team cannot be the same")
	}

	if req.Season != nil {
		game.Season = *req.Season
	}

	if req.Week != nil {
		game.Week = *req.Week
	}

	if req.GameDate != nil {
		game.GameDate = *req.GameDate
	}

	if req.Status != nil {
		game.Status = *req.Status
	}

	if req.HomeScore != nil {
		game.HomeScore = req.HomeScore
	}

	if req.AwayScore != nil {
		game.AwayScore = req.AwayScore
	}

	// Update the game
	if err := s.gameRepo.Update(game); err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	return game, nil
}

// DeleteGame deletes a game by ID
func (s *gameService) DeleteGame(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid game ID: %d", id)
	}

	// Check if game exists
	exists, err := s.gameRepo.Exists(id)
	if err != nil {
		return fmt.Errorf("failed to check if game exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("game with ID %d not found", id)
	}

	return s.gameRepo.Delete(id)
}

// GetGamesByTeam retrieves all games for a specific team
func (s *gameService) GetGamesByTeam(teamID int) ([]*models.Game, error) {
	if teamID <= 0 {
		return nil, fmt.Errorf("invalid team ID: %d", teamID)
	}

	// Check if team exists
	exists, err := s.teamRepo.Exists(teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if team exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("team with ID %d not found", teamID)
	}

	return s.gameRepo.GetByTeamID(teamID)
}

// GetGamesBySeason retrieves all games for a specific season
func (s *gameService) GetGamesBySeason(season string) ([]*models.Game, error) {
	if season == "" {
		return nil, fmt.Errorf("season cannot be empty")
	}

	return s.gameRepo.GetBySeason(season)
}

// GetGamesByWeek retrieves all games for a specific week in a season
func (s *gameService) GetGamesByWeek(season string, week int) ([]*models.Game, error) {
	if season == "" {
		return nil, fmt.Errorf("season cannot be empty")
	}

	if week < 1 || week > 22 {
		return nil, fmt.Errorf("week must be between 1 and 22, got %d", week)
	}

	return s.gameRepo.GetByWeek(season, week)
}

// validateCreateGameRequest validates a create game request
func (s *gameService) validateCreateGameRequest(req *models.CreateGameRequest) error {
	if req.HomeTeamID <= 0 {
		return fmt.Errorf("home team ID must be positive")
	}

	if req.AwayTeamID <= 0 {
		return fmt.Errorf("away team ID must be positive")
	}

	if req.Season == "" {
		return fmt.Errorf("season is required")
	}

	if req.Week < 1 || req.Week > 22 {
		return fmt.Errorf("week must be between 1 and 22, got %d", req.Week)
	}

	if req.GameDate.IsZero() {
		return fmt.Errorf("game date is required")
	}

	// Check if game date is not too far in the past (more than 1 year)
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	if req.GameDate.Before(oneYearAgo) {
		return fmt.Errorf("game date cannot be more than 1 year in the past")
	}

	// Check if game date is not too far in the future (more than 2 years)
	twoYearsFromNow := time.Now().AddDate(2, 0, 0)
	if req.GameDate.After(twoYearsFromNow) {
		return fmt.Errorf("game date cannot be more than 2 years in the future")
	}

	if req.Status != "" {
		validStatuses := []string{"scheduled", "in_progress", "completed", "cancelled"}
		valid := false
		for _, status := range validStatuses {
			if req.Status == status {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid status: %s. Must be one of: scheduled, in_progress, completed, cancelled", req.Status)
		}
	}

	if req.HomeScore != nil && *req.HomeScore < 0 {
		return fmt.Errorf("home score cannot be negative")
	}

	if req.AwayScore != nil && *req.AwayScore < 0 {
		return fmt.Errorf("away score cannot be negative")
	}

	return nil
}

// validateUpdateGameRequest validates an update game request
func (s *gameService) validateUpdateGameRequest(req *models.UpdateGameRequest) error {
	if req.HomeTeamID != nil && *req.HomeTeamID <= 0 {
		return fmt.Errorf("home team ID must be positive")
	}

	if req.AwayTeamID != nil && *req.AwayTeamID <= 0 {
		return fmt.Errorf("away team ID must be positive")
	}

	if req.Season != nil && *req.Season == "" {
		return fmt.Errorf("season cannot be empty")
	}

	if req.Week != nil && (*req.Week < 1 || *req.Week > 22) {
		return fmt.Errorf("week must be between 1 and 22, got %d", *req.Week)
	}

	if req.GameDate != nil {
		if req.GameDate.IsZero() {
			return fmt.Errorf("game date cannot be zero")
		}

		// Check if game date is not too far in the past (more than 1 year)
		oneYearAgo := time.Now().AddDate(-1, 0, 0)
		if req.GameDate.Before(oneYearAgo) {
			return fmt.Errorf("game date cannot be more than 1 year in the past")
		}

		// Check if game date is not too far in the future (more than 2 years)
		twoYearsFromNow := time.Now().AddDate(2, 0, 0)
		if req.GameDate.After(twoYearsFromNow) {
			return fmt.Errorf("game date cannot be more than 2 years in the future")
		}
	}

	if req.Status != nil {
		validStatuses := []string{"scheduled", "in_progress", "completed", "cancelled"}
		valid := false
		for _, status := range validStatuses {
			if *req.Status == status {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid status: %s. Must be one of: scheduled, in_progress, completed, cancelled", *req.Status)
		}
	}

	if req.HomeScore != nil && *req.HomeScore < 0 {
		return fmt.Errorf("home score cannot be negative")
	}

	if req.AwayScore != nil && *req.AwayScore < 0 {
		return fmt.Errorf("away score cannot be negative")
	}

	return nil
}
