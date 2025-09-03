package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sports-backend/models"
	"sports-backend/services"

	"github.com/gorilla/mux"
)

// PlayerHandler handles HTTP requests for players
type PlayerHandler struct {
	playerService services.PlayerService
}

// NewPlayerHandler creates a new player handler
func NewPlayerHandler(playerService services.PlayerService) *PlayerHandler {
	return &PlayerHandler{
		playerService: playerService,
	}
}

// GetPlayers handles GET /api/players
func (h *PlayerHandler) GetPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := h.playerService.GetAllPlayers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}

// CreatePlayer handles POST /api/players
func (h *PlayerHandler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	player, err := h.playerService.CreatePlayer(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(player)
}

// GetPlayer handles GET /api/players/{id}
func (h *PlayerHandler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid player ID", http.StatusBadRequest)
		return
	}

	player, err := h.playerService.GetPlayer(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

// UpdatePlayer handles PUT /api/players/{id}
func (h *PlayerHandler) UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid player ID", http.StatusBadRequest)
		return
	}

	var req models.UpdatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	player, err := h.playerService.UpdatePlayer(id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

// DeletePlayer handles DELETE /api/players/{id}
func (h *PlayerHandler) DeletePlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid player ID", http.StatusBadRequest)
		return
	}

	if err := h.playerService.DeletePlayer(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetPlayerStats handles GET /api/players/{id}/stats
func (h *PlayerHandler) GetPlayerStats(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement when player stats service is created
	http.Error(w, "Not implemented yet", http.StatusNotImplemented)
}

// CreatePlayerStats handles POST /api/players/{id}/stats
func (h *PlayerHandler) CreatePlayerStats(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement when player stats service is created
	http.Error(w, "Not implemented yet", http.StatusNotImplemented)
}

// DeletePlayerStats handles DELETE /api/players/{id}/stats/{stats_id}
func (h *PlayerHandler) DeletePlayerStats(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement when player stats service is created
	http.Error(w, "Not implemented yet", http.StatusNotImplemented)
}

// UpdatePlayerStats handles PUT /api/players/{id}/stats/{stats_id}
func (h *PlayerHandler) UpdatePlayerStats(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement when player stats service is created
	http.Error(w, "Not implemented yet", http.StatusNotImplemented)
}
