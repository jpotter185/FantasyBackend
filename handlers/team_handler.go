package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sports-backend/models"
	"sports-backend/services"

	"github.com/gorilla/mux"
)

// TeamHandler handles HTTP requests for teams
type TeamHandler struct {
	teamService services.TeamService
}

// NewTeamHandler creates a new team handler
func NewTeamHandler(teamService services.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

// GetTeams handles GET /api/teams
func (h *TeamHandler) GetTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.teamService.GetAllTeams()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

// CreateTeam handles POST /api/teams
func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	team, err := h.teamService.CreateTeam(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(team)
}

// GetTeam handles GET /api/teams/{id}
func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	team, err := h.teamService.GetTeam(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

// UpdateTeam handles PUT /api/teams/{id}
func (h *TeamHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	team, err := h.teamService.UpdateTeam(id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

// DeleteTeam handles DELETE /api/teams/{id}
func (h *TeamHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	if err := h.teamService.DeleteTeam(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetTeamStats handles GET /api/teams/{id}/stats
func (h *TeamHandler) GetTeamStats(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement when team stats service is created
	http.Error(w, "Not implemented yet", http.StatusNotImplemented)
}

// CreateTeamStats handles POST /api/teams/{id}/stats
func (h *TeamHandler) CreateTeamStats(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement when team stats service is created
	http.Error(w, "Not implemented yet", http.StatusNotImplemented)
}
