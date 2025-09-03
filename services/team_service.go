package services

import (
	"fmt"
	"strings"

	"sports-backend/models"
	"sports-backend/repositories"
)

// TeamService defines the interface for team business logic
type TeamService interface {
	GetTeam(id int) (*models.Team, error)
	GetAllTeams() ([]*models.Team, error)
	GetTeamsByConference(conference string) ([]*models.Team, error)
	GetTeamsByDivision(division string) ([]*models.Team, error)
	CreateTeam(req *models.CreateTeamRequest) (*models.Team, error)
	UpdateTeam(id int, req *models.UpdateTeamRequest) (*models.Team, error)
	DeleteTeam(id int) error
}

// teamService implements TeamService interface
type teamService struct {
	teamRepo repositories.TeamRepository
}

// NewTeamService creates a new team service
func NewTeamService(teamRepo repositories.TeamRepository) TeamService {
	return &teamService{
		teamRepo: teamRepo,
	}
}

// GetTeam retrieves a team by ID
func (s *teamService) GetTeam(id int) (*models.Team, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid team ID: %d", id)
	}

	team, err := s.teamRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return team, nil
}

// GetAllTeams retrieves all teams
func (s *teamService) GetAllTeams() ([]*models.Team, error) {
	teams, err := s.teamRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get teams: %w", err)
	}

	return teams, nil
}

// GetTeamsByConference retrieves all teams in a specific conference
func (s *teamService) GetTeamsByConference(conference string) ([]*models.Team, error) {
	if strings.TrimSpace(conference) == "" {
		return nil, fmt.Errorf("conference cannot be empty")
	}

	teams, err := s.teamRepo.GetByConference(strings.TrimSpace(conference))
	if err != nil {
		return nil, fmt.Errorf("failed to get teams by conference: %w", err)
	}

	return teams, nil
}

// GetTeamsByDivision retrieves all teams in a specific division
func (s *teamService) GetTeamsByDivision(division string) ([]*models.Team, error) {
	if strings.TrimSpace(division) == "" {
		return nil, fmt.Errorf("division cannot be empty")
	}

	teams, err := s.teamRepo.GetByDivision(strings.TrimSpace(division))
	if err != nil {
		return nil, fmt.Errorf("failed to get teams by division: %w", err)
	}

	return teams, nil
}

// CreateTeam creates a new team
func (s *teamService) CreateTeam(req *models.CreateTeamRequest) (*models.Team, error) {
	// Validate request
	if err := s.validateCreateTeamRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create team
	team := &models.Team{
		Name:       strings.TrimSpace(req.Name),
		City:       strings.TrimSpace(req.City),
		Conference: strings.TrimSpace(req.Conference),
		Division:   strings.TrimSpace(req.Division),
	}

	if err := s.teamRepo.Create(team); err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	return team, nil
}

// UpdateTeam updates an existing team
func (s *teamService) UpdateTeam(id int, req *models.UpdateTeamRequest) (*models.Team, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid team ID: %d", id)
	}

	// Validate request
	if err := s.validateUpdateTeamRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing team
	team, err := s.teamRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		team.Name = strings.TrimSpace(*req.Name)
	}
	if req.City != nil {
		team.City = strings.TrimSpace(*req.City)
	}
	if req.Conference != nil {
		team.Conference = strings.TrimSpace(*req.Conference)
	}
	if req.Division != nil {
		team.Division = strings.TrimSpace(*req.Division)
	}

	// Update team
	if err := s.teamRepo.Update(team); err != nil {
		return nil, fmt.Errorf("failed to update team: %w", err)
	}

	return team, nil
}

// DeleteTeam deletes a team
func (s *teamService) DeleteTeam(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid team ID: %d", id)
	}

	// Check if team exists
	exists, err := s.teamRepo.Exists(id)
	if err != nil {
		return fmt.Errorf("failed to check team existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("team with ID %d not found", id)
	}

	// TODO: Add business logic here if needed
	// For example: check if team has players, prevent deletion if they do
	// For example: check if team has games, prevent deletion if they do

	if err := s.teamRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	return nil
}

// validateCreateTeamRequest validates the create team request
func (s *teamService) validateCreateTeamRequest(req *models.CreateTeamRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("team name is required")
	}

	if strings.TrimSpace(req.City) == "" {
		return fmt.Errorf("city is required")
	}

	if strings.TrimSpace(req.Conference) == "" {
		return fmt.Errorf("conference is required")
	}

	if strings.TrimSpace(req.Division) == "" {
		return fmt.Errorf("division is required")
	}

	// Validate conference and division values
	validConferences := []string{"AFC", "NFC"}
	validDivisions := []string{"North", "South", "East", "West"}

	conferenceValid := false
	for _, validConf := range validConferences {
		if strings.EqualFold(req.Conference, validConf) {
			conferenceValid = true
			break
		}
	}
	if !conferenceValid {
		return fmt.Errorf("conference must be one of: %v", validConferences)
	}

	divisionValid := false
	for _, validDiv := range validDivisions {
		if strings.EqualFold(req.Division, validDiv) {
			divisionValid = true
			break
		}
	}
	if !divisionValid {
		return fmt.Errorf("division must be one of: %v", validDivisions)
	}

	return nil
}

// validateUpdateTeamRequest validates the update team request
func (s *teamService) validateUpdateTeamRequest(req *models.UpdateTeamRequest) error {
	// Check if at least one field is being updated
	if req.Name == nil && req.City == nil && req.Conference == nil &&
		req.Division == nil {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Validate individual fields if provided
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return fmt.Errorf("team name cannot be empty")
	}

	if req.City != nil && strings.TrimSpace(*req.City) == "" {
		return fmt.Errorf("city cannot be empty")
	}

	if req.Conference != nil && strings.TrimSpace(*req.Conference) == "" {
		return fmt.Errorf("conference cannot be empty")
	}

	if req.Division != nil && strings.TrimSpace(*req.Division) == "" {
		return fmt.Errorf("division cannot be empty")
	}

	// Validate conference and division values if provided
	if req.Conference != nil {
		validConferences := []string{"AFC", "NFC"}
		conferenceValid := false
		for _, validConf := range validConferences {
			if strings.EqualFold(*req.Conference, validConf) {
				conferenceValid = true
				break
			}
		}
		if !conferenceValid {
			return fmt.Errorf("conference must be one of: %v", validConferences)
		}
	}

	if req.Division != nil {
		validDivisions := []string{"North", "South", "East", "West"}
		divisionValid := false
		for _, validDiv := range validDivisions {
			if strings.EqualFold(*req.Division, validDiv) {
				divisionValid = true
				break
			}
		}
		if !divisionValid {
			return fmt.Errorf("division must be one of: %v", validDivisions)
		}
	}

	return nil
}
