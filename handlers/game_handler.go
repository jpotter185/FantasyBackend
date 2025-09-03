package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sports-backend/models"
	"sports-backend/services"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// GameHandler handles HTTP requests for games
type GameHandler struct {
	gameService services.GameService
}

// NewGameHandler creates a new game handler
func NewGameHandler(gameService services.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

// GetGames handles GET /api/games
func (h *GameHandler) GetGames(w http.ResponseWriter, r *http.Request) {
	games, err := h.gameService.GetAllGames()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get games: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

// GetGame handles GET /api/games/{id}
func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	game, err := h.gameService.GetGameByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to get game: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

// CreateGame handles POST /api/games
func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var req models.CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	game, err := h.gameService.CreateGame(&req)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") ||
			strings.Contains(err.Error(), "not found") ||
			strings.Contains(err.Error(), "cannot be the same") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to create game: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(game)
}

// UpdateGame handles PUT /api/games/{id}
func (h *GameHandler) UpdateGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	game, err := h.gameService.UpdateGame(id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") ||
			strings.Contains(err.Error(), "not found") ||
			strings.Contains(err.Error(), "cannot be the same") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to update game: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

// DeleteGame handles DELETE /api/games/{id}
func (h *GameHandler) DeleteGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	err = h.gameService.DeleteGame(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to delete game: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetGamesByTeam handles GET /api/teams/{id}/games
func (h *GameHandler) GetGamesByTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamIDStr := vars["id"]

	teamID, err := strconv.Atoi(teamIDStr)
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	games, err := h.gameService.GetGamesByTeam(teamID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to get games: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

// GetGamesBySeason handles GET /api/games/season/{season}
func (h *GameHandler) GetGamesBySeason(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	season := vars["season"]

	if season == "" {
		http.Error(w, "Season parameter is required", http.StatusBadRequest)
		return
	}

	games, err := h.gameService.GetGamesBySeason(season)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get games: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

// GetGamesByWeek handles GET /api/games/season/{season}/week/{week}
func (h *GameHandler) GetGamesByWeek(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	season := vars["season"]
	weekStr := vars["week"]

	if season == "" {
		http.Error(w, "Season parameter is required", http.StatusBadRequest)
		return
	}

	week, err := strconv.Atoi(weekStr)
	if err != nil {
		http.Error(w, "Invalid week parameter", http.StatusBadRequest)
		return
	}

	games, err := h.gameService.GetGamesByWeek(season, week)
	if err != nil {
		if strings.Contains(err.Error(), "must be between 1 and 22") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to get games: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}
