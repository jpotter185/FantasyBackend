package services

import (
	"fmt"

	"sports-backend/models"
	"sports-backend/repositories"
)

// PlayerStatsService defines the interface for player stats business logic
type PlayerStatsService interface {
	GetPlayerStats(id int) (*models.PlayerStats, error)
	GetAllPlayerStats() ([]*models.PlayerStats, error)
	GetPlayerStatsByPlayer(playerID int) ([]*models.PlayerStats, error)
	GetPlayerStatsByGame(gameID int) ([]*models.PlayerStats, error)
	CreatePlayerStats(req *models.CreatePlayerStatsRequest) (*models.PlayerStats, error)
	UpdatePlayerStats(id int, req *models.UpdatePlayerStatsRequest) (*models.PlayerStats, error)
	DeletePlayerStats(id int) error
}

// playerStatsService implements PlayerStatsService interface
type playerStatsService struct {
	playerStatsRepo repositories.PlayerStatsRepository
	playerRepo      repositories.PlayerRepository
}

// NewPlayerStatsService creates a new player stats service
func NewPlayerStatsService(playerStatsRepo repositories.PlayerStatsRepository, playerRepo repositories.PlayerRepository) PlayerStatsService {
	return &playerStatsService{
		playerStatsRepo: playerStatsRepo,
		playerRepo:      playerRepo,
	}
}

// GetPlayerStats retrieves player stats by ID
func (s *playerStatsService) GetPlayerStats(id int) (*models.PlayerStats, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid player stats ID: %d", id)
	}

	stats, err := s.playerStatsRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get player stats: %w", err)
	}

	return stats, nil
}

// GetAllPlayerStats retrieves all player stats
func (s *playerStatsService) GetAllPlayerStats() ([]*models.PlayerStats, error) {
	statsList, err := s.playerStatsRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all player stats: %w", err)
	}

	return statsList, nil
}

// GetPlayerStatsByPlayer retrieves all stats for a specific player
func (s *playerStatsService) GetPlayerStatsByPlayer(playerID int) ([]*models.PlayerStats, error) {
	if playerID <= 0 {
		return nil, fmt.Errorf("invalid player ID: %d", playerID)
	}

	// Verify player exists
	exists, err := s.playerRepo.Exists(playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify player existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("player with ID %d not found", playerID)
	}

	statsList, err := s.playerStatsRepo.GetByPlayerID(playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player stats by player: %w", err)
	}

	return statsList, nil
}

// GetPlayerStatsByGame retrieves all stats for a specific game
func (s *playerStatsService) GetPlayerStatsByGame(gameID int) ([]*models.PlayerStats, error) {
	if gameID <= 0 {
		return nil, fmt.Errorf("invalid game ID: %d", gameID)
	}

	// TODO: Verify game exists when game repository is implemented
	// For now, we'll skip this validation

	statsList, err := s.playerStatsRepo.GetByGameID(gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player stats by game: %w", err)
	}

	return statsList, nil
}

// CreatePlayerStats creates new player stats
func (s *playerStatsService) CreatePlayerStats(req *models.CreatePlayerStatsRequest) (*models.PlayerStats, error) {
	// Validate request
	if err := s.validateCreatePlayerStatsRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Verify player exists
	exists, err := s.playerRepo.Exists(req.PlayerID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify player existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("player with ID %d not found", req.PlayerID)
	}

	// Check if stats already exist for this player and game
	exists, err = s.playerStatsRepo.ExistsByPlayerAndGame(req.PlayerID, req.GameID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing stats: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("player stats already exist for player %d in game %d", req.PlayerID, req.GameID)
	}

	// Create player stats
	stats := &models.PlayerStats{
		PlayerID:               req.PlayerID,
		GameID:                 req.GameID,
		PassingAttempts:        req.PassingAttempts,
		PassingCompletions:     req.PassingCompletions,
		PassingYards:           req.PassingYards,
		PassingTouchdowns:      req.PassingTouchdowns,
		PassingInterceptions:   req.PassingInterceptions,
		RushingAttempts:        req.RushingAttempts,
		RushingYards:           req.RushingYards,
		RushingTouchdowns:      req.RushingTouchdowns,
		ReceivingTargets:       req.ReceivingTargets,
		Receptions:             req.Receptions,
		ReceivingYards:         req.ReceivingYards,
		ReceivingTouchdowns:    req.ReceivingTouchdowns,
		Fumbles:                req.Fumbles,
		FumblesLost:            req.FumblesLost,
		Tackles:                req.Tackles,
		SoloTackles:            req.SoloTackles,
		AssistedTackles:        req.AssistedTackles,
		Sacks:                  req.Sacks,
		DefensiveInterceptions: req.DefensiveInterceptions,
		PassDeflections:        req.PassDeflections,
		ForcedFumbles:          req.ForcedFumbles,
		FumbleRecoveries:       req.FumbleRecoveries,
		DefensiveTouchdowns:    req.DefensiveTouchdowns,
		FieldGoalsAttempted:    req.FieldGoalsAttempted,
		FieldGoalsMade:         req.FieldGoalsMade,
		ExtraPointsAttempted:   req.ExtraPointsAttempted,
		ExtraPointsMade:        req.ExtraPointsMade,
		Punts:                  req.Punts,
		PuntYards:              req.PuntYards,
		KickReturns:            req.KickReturns,
		KickReturnYards:        req.KickReturnYards,
		KickReturnTouchdowns:   req.KickReturnTouchdowns,
		PuntReturns:            req.PuntReturns,
		PuntReturnYards:        req.PuntReturnYards,
		PuntReturnTouchdowns:   req.PuntReturnTouchdowns,
	}

	if err := s.playerStatsRepo.Create(stats); err != nil {
		return nil, fmt.Errorf("failed to create player stats: %w", err)
	}

	return stats, nil
}

// UpdatePlayerStats updates existing player stats
func (s *playerStatsService) UpdatePlayerStats(id int, req *models.UpdatePlayerStatsRequest) (*models.PlayerStats, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid player stats ID: %d", id)
	}

	// Validate request
	if err := s.validateUpdatePlayerStatsRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing stats
	stats, err := s.playerStatsRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get player stats: %w", err)
	}

	// Update fields if provided
	if req.PassingAttempts != nil {
		stats.PassingAttempts = req.PassingAttempts
	}
	if req.PassingCompletions != nil {
		stats.PassingCompletions = req.PassingCompletions
	}
	if req.PassingYards != nil {
		stats.PassingYards = req.PassingYards
	}
	if req.PassingTouchdowns != nil {
		stats.PassingTouchdowns = req.PassingTouchdowns
	}
	if req.PassingInterceptions != nil {
		stats.PassingInterceptions = req.PassingInterceptions
	}
	if req.RushingAttempts != nil {
		stats.RushingAttempts = req.RushingAttempts
	}
	if req.RushingYards != nil {
		stats.RushingYards = req.RushingYards
	}
	if req.RushingTouchdowns != nil {
		stats.RushingTouchdowns = req.RushingTouchdowns
	}
	if req.ReceivingTargets != nil {
		stats.ReceivingTargets = req.ReceivingTargets
	}
	if req.Receptions != nil {
		stats.Receptions = req.Receptions
	}
	if req.ReceivingYards != nil {
		stats.ReceivingYards = req.ReceivingYards
	}
	if req.ReceivingTouchdowns != nil {
		stats.ReceivingTouchdowns = req.ReceivingTouchdowns
	}
	if req.Fumbles != nil {
		stats.Fumbles = req.Fumbles
	}
	if req.FumblesLost != nil {
		stats.FumblesLost = req.FumblesLost
	}
	if req.Tackles != nil {
		stats.Tackles = req.Tackles
	}
	if req.SoloTackles != nil {
		stats.SoloTackles = req.SoloTackles
	}
	if req.AssistedTackles != nil {
		stats.AssistedTackles = req.AssistedTackles
	}
	if req.Sacks != nil {
		stats.Sacks = req.Sacks
	}
	if req.DefensiveInterceptions != nil {
		stats.DefensiveInterceptions = req.DefensiveInterceptions
	}
	if req.PassDeflections != nil {
		stats.PassDeflections = req.PassDeflections
	}
	if req.ForcedFumbles != nil {
		stats.ForcedFumbles = req.ForcedFumbles
	}
	if req.FumbleRecoveries != nil {
		stats.FumbleRecoveries = req.FumbleRecoveries
	}
	if req.DefensiveTouchdowns != nil {
		stats.DefensiveTouchdowns = req.DefensiveTouchdowns
	}
	if req.FieldGoalsAttempted != nil {
		stats.FieldGoalsAttempted = req.FieldGoalsAttempted
	}
	if req.FieldGoalsMade != nil {
		stats.FieldGoalsMade = req.FieldGoalsMade
	}
	if req.ExtraPointsAttempted != nil {
		stats.ExtraPointsAttempted = req.ExtraPointsAttempted
	}
	if req.ExtraPointsMade != nil {
		stats.ExtraPointsMade = req.ExtraPointsMade
	}
	if req.Punts != nil {
		stats.Punts = req.Punts
	}
	if req.PuntYards != nil {
		stats.PuntYards = req.PuntYards
	}
	if req.KickReturns != nil {
		stats.KickReturns = req.KickReturns
	}
	if req.KickReturnYards != nil {
		stats.KickReturnYards = req.KickReturnYards
	}
	if req.KickReturnTouchdowns != nil {
		stats.KickReturnTouchdowns = req.KickReturnTouchdowns
	}
	if req.PuntReturns != nil {
		stats.PuntReturns = req.PuntReturns
	}
	if req.PuntReturnYards != nil {
		stats.PuntReturnYards = req.PuntReturnYards
	}
	if req.PuntReturnTouchdowns != nil {
		stats.PuntReturnTouchdowns = req.PuntReturnTouchdowns
	}

	// Update stats
	if err := s.playerStatsRepo.Update(stats); err != nil {
		return nil, fmt.Errorf("failed to update player stats: %w", err)
	}

	return stats, nil
}

// DeletePlayerStats deletes player stats
func (s *playerStatsService) DeletePlayerStats(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid player stats ID: %d", id)
	}

	// Check if stats exist
	exists, err := s.playerStatsRepo.Exists(id)
	if err != nil {
		return fmt.Errorf("failed to check player stats existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("player stats with ID %d not found", id)
	}

	if err := s.playerStatsRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete player stats: %w", err)
	}

	return nil
}

// validateCreatePlayerStatsRequest validates the create player stats request
func (s *playerStatsService) validateCreatePlayerStatsRequest(req *models.CreatePlayerStatsRequest) error {
	if req.PlayerID <= 0 {
		return fmt.Errorf("player ID is required and must be positive")
	}

	if req.GameID <= 0 {
		return fmt.Errorf("game ID is required and must be positive")
	}

	// Validate that at least one stat is provided
	if req.PassingAttempts == nil && req.PassingCompletions == nil && req.PassingYards == nil &&
		req.PassingTouchdowns == nil && req.PassingInterceptions == nil &&
		req.RushingAttempts == nil && req.RushingYards == nil && req.RushingTouchdowns == nil &&
		req.ReceivingTargets == nil && req.Receptions == nil && req.ReceivingYards == nil &&
		req.ReceivingTouchdowns == nil && req.Fumbles == nil && req.FumblesLost == nil &&
		req.Tackles == nil && req.SoloTackles == nil && req.AssistedTackles == nil &&
		req.Sacks == nil && req.DefensiveInterceptions == nil && req.PassDeflections == nil &&
		req.ForcedFumbles == nil && req.FumbleRecoveries == nil && req.DefensiveTouchdowns == nil &&
		req.FieldGoalsAttempted == nil && req.FieldGoalsMade == nil &&
		req.ExtraPointsAttempted == nil && req.ExtraPointsMade == nil &&
		req.Punts == nil && req.PuntYards == nil && req.KickReturns == nil &&
		req.KickReturnYards == nil && req.KickReturnTouchdowns == nil &&
		req.PuntReturns == nil && req.PuntReturnYards == nil && req.PuntReturnTouchdowns == nil {
		return fmt.Errorf("at least one statistic must be provided")
	}

	// Validate logical constraints
	if err := s.validateStatConstraints(req); err != nil {
		return err
	}

	return nil
}

// validateUpdatePlayerStatsRequest validates the update player stats request
func (s *playerStatsService) validateUpdatePlayerStatsRequest(req *models.UpdatePlayerStatsRequest) error {
	// Check if at least one field is being updated
	if req.PassingAttempts == nil && req.PassingCompletions == nil && req.PassingYards == nil &&
		req.PassingTouchdowns == nil && req.PassingInterceptions == nil &&
		req.RushingAttempts == nil && req.RushingYards == nil && req.RushingTouchdowns == nil &&
		req.ReceivingTargets == nil && req.Receptions == nil && req.ReceivingYards == nil &&
		req.ReceivingTouchdowns == nil && req.Fumbles == nil && req.FumblesLost == nil &&
		req.Tackles == nil && req.SoloTackles == nil && req.AssistedTackles == nil &&
		req.Sacks == nil && req.DefensiveInterceptions == nil && req.PassDeflections == nil &&
		req.ForcedFumbles == nil && req.FumbleRecoveries == nil && req.DefensiveTouchdowns == nil &&
		req.FieldGoalsAttempted == nil && req.FieldGoalsMade == nil &&
		req.ExtraPointsAttempted == nil && req.ExtraPointsMade == nil &&
		req.Punts == nil && req.PuntYards == nil && req.KickReturns == nil &&
		req.KickReturnYards == nil && req.KickReturnTouchdowns == nil &&
		req.PuntReturns == nil && req.PuntReturnYards == nil && req.PuntReturnTouchdowns == nil {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Validate logical constraints
	if err := s.validateUpdateStatConstraints(req); err != nil {
		return err
	}

	return nil
}

// validateStatConstraints validates logical constraints for create requests
func (s *playerStatsService) validateStatConstraints(req *models.CreatePlayerStatsRequest) error {
	// Passing completions cannot exceed passing attempts
	if req.PassingCompletions != nil && req.PassingAttempts != nil {
		if *req.PassingCompletions > *req.PassingAttempts {
			return fmt.Errorf("passing completions cannot exceed passing attempts")
		}
	}

	// Solo tackles + assisted tackles should equal total tackles (if all provided)
	if req.Tackles != nil && req.SoloTackles != nil && req.AssistedTackles != nil {
		if *req.Tackles != *req.SoloTackles+*req.AssistedTackles {
			return fmt.Errorf("total tackles must equal solo tackles plus assisted tackles")
		}
	}

	// Field goals made cannot exceed field goals attempted
	if req.FieldGoalsMade != nil && req.FieldGoalsAttempted != nil {
		if *req.FieldGoalsMade > *req.FieldGoalsAttempted {
			return fmt.Errorf("field goals made cannot exceed field goals attempted")
		}
	}

	// Extra points made cannot exceed extra points attempted
	if req.ExtraPointsMade != nil && req.ExtraPointsAttempted != nil {
		if *req.ExtraPointsMade > *req.ExtraPointsAttempted {
			return fmt.Errorf("extra points made cannot exceed extra points attempted")
		}
	}

	// Fumbles lost cannot exceed total fumbles
	if req.FumblesLost != nil && req.Fumbles != nil {
		if *req.FumblesLost > *req.Fumbles {
			return fmt.Errorf("fumbles lost cannot exceed total fumbles")
		}
	}

	// Validate non-negative values
	nonNegativeFields := []struct {
		value *int
		name  string
	}{
		{req.PassingAttempts, "passing attempts"},
		{req.PassingCompletions, "passing completions"},
		{req.PassingYards, "passing yards"},
		{req.PassingTouchdowns, "passing touchdowns"},
		{req.PassingInterceptions, "passing interceptions"},
		{req.RushingAttempts, "rushing attempts"},
		{req.RushingYards, "rushing yards"},
		{req.RushingTouchdowns, "rushing touchdowns"},
		{req.ReceivingTargets, "receiving targets"},
		{req.Receptions, "receptions"},
		{req.ReceivingYards, "receiving yards"},
		{req.ReceivingTouchdowns, "receiving touchdowns"},
		{req.Fumbles, "fumbles"},
		{req.FumblesLost, "fumbles lost"},
		{req.Tackles, "tackles"},
		{req.SoloTackles, "solo tackles"},
		{req.AssistedTackles, "assisted tackles"},
		{req.Sacks, "sacks"},
		{req.DefensiveInterceptions, "defensive interceptions"},
		{req.PassDeflections, "pass deflections"},
		{req.ForcedFumbles, "forced fumbles"},
		{req.FumbleRecoveries, "fumble recoveries"},
		{req.DefensiveTouchdowns, "defensive touchdowns"},
		{req.FieldGoalsAttempted, "field goals attempted"},
		{req.FieldGoalsMade, "field goals made"},
		{req.ExtraPointsAttempted, "extra points attempted"},
		{req.ExtraPointsMade, "extra points made"},
		{req.Punts, "punts"},
		{req.PuntYards, "punt yards"},
		{req.KickReturns, "kick returns"},
		{req.KickReturnYards, "kick return yards"},
		{req.KickReturnTouchdowns, "kick return touchdowns"},
		{req.PuntReturns, "punt returns"},
		{req.PuntReturnYards, "punt return yards"},
		{req.PuntReturnTouchdowns, "punt return touchdowns"},
	}

	for _, field := range nonNegativeFields {
		if field.value != nil && *field.value < 0 {
			return fmt.Errorf("%s cannot be negative", field.name)
		}
	}

	return nil
}

// validateUpdateStatConstraints validates logical constraints for update requests
func (s *playerStatsService) validateUpdateStatConstraints(req *models.UpdatePlayerStatsRequest) error {
	// Passing completions cannot exceed passing attempts
	if req.PassingCompletions != nil && req.PassingAttempts != nil {
		if *req.PassingCompletions > *req.PassingAttempts {
			return fmt.Errorf("passing completions cannot exceed passing attempts")
		}
	}

	// Solo tackles + assisted tackles should equal total tackles (if all provided)
	if req.Tackles != nil && req.SoloTackles != nil && req.AssistedTackles != nil {
		if *req.Tackles != *req.SoloTackles+*req.AssistedTackles {
			return fmt.Errorf("total tackles must equal solo tackles plus assisted tackles")
		}
	}

	// Field goals made cannot exceed field goals attempted
	if req.FieldGoalsMade != nil && req.FieldGoalsAttempted != nil {
		if *req.FieldGoalsMade > *req.FieldGoalsAttempted {
			return fmt.Errorf("field goals made cannot exceed field goals attempted")
		}
	}

	// Extra points made cannot exceed extra points attempted
	if req.ExtraPointsMade != nil && req.ExtraPointsAttempted != nil {
		if *req.ExtraPointsMade > *req.ExtraPointsAttempted {
			return fmt.Errorf("extra points made cannot exceed extra points attempted")
		}
	}

	// Fumbles lost cannot exceed total fumbles
	if req.FumblesLost != nil && req.Fumbles != nil {
		if *req.FumblesLost > *req.Fumbles {
			return fmt.Errorf("fumbles lost cannot exceed total fumbles")
		}
	}

	// Validate non-negative values
	nonNegativeFields := []struct {
		value *int
		name  string
	}{
		{req.PassingAttempts, "passing attempts"},
		{req.PassingCompletions, "passing completions"},
		{req.PassingYards, "passing yards"},
		{req.PassingTouchdowns, "passing touchdowns"},
		{req.PassingInterceptions, "passing interceptions"},
		{req.RushingAttempts, "rushing attempts"},
		{req.RushingYards, "rushing yards"},
		{req.RushingTouchdowns, "rushing touchdowns"},
		{req.ReceivingTargets, "receiving targets"},
		{req.Receptions, "receptions"},
		{req.ReceivingYards, "receiving yards"},
		{req.ReceivingTouchdowns, "receiving touchdowns"},
		{req.Fumbles, "fumbles"},
		{req.FumblesLost, "fumbles lost"},
		{req.Tackles, "tackles"},
		{req.SoloTackles, "solo tackles"},
		{req.AssistedTackles, "assisted tackles"},
		{req.Sacks, "sacks"},
		{req.DefensiveInterceptions, "defensive interceptions"},
		{req.PassDeflections, "pass deflections"},
		{req.ForcedFumbles, "forced fumbles"},
		{req.FumbleRecoveries, "fumble recoveries"},
		{req.DefensiveTouchdowns, "defensive touchdowns"},
		{req.FieldGoalsAttempted, "field goals attempted"},
		{req.FieldGoalsMade, "field goals made"},
		{req.ExtraPointsAttempted, "extra points attempted"},
		{req.ExtraPointsMade, "extra points made"},
		{req.Punts, "punts"},
		{req.PuntYards, "punt yards"},
		{req.KickReturns, "kick returns"},
		{req.KickReturnYards, "kick return yards"},
		{req.KickReturnTouchdowns, "kick return touchdowns"},
		{req.PuntReturns, "punt returns"},
		{req.PuntReturnYards, "punt return yards"},
		{req.PuntReturnTouchdowns, "punt return touchdowns"},
	}

	for _, field := range nonNegativeFields {
		if field.value != nil && *field.value < 0 {
			return fmt.Errorf("%s cannot be negative", field.name)
		}
	}

	return nil
}
