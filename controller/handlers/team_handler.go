package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/gonzabosio/res-manager/model"
)

func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()
	team := new(model.Team)
	json.NewDecoder(r.Body).Decode(&team)
	err := validate.Struct(team)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		WriteJSON(w, map[string]string{
			"message": fmt.Sprintf("Validation error: %s", errors),
		}, http.StatusBadRequest)
		return
	}
	id, err := h.RepositoryService.CreateTeam(team)
	if err != nil {
		WriteJSON(w, map[string]string{
			"message": err.Error(),
		}, http.StatusBadRequest)
		return
	}
	WriteJSON(w, map[string]string{
		"message": "Team created successfully",
		"team_id": id,
	}, http.StatusOK)
}

func (h *Handler) GetTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.RepositoryService.ReadTeams()
	if err != nil {
		WriteJSON(w, map[string]string{
			"error": err.Error(),
		}, http.StatusBadRequest)
		return
	}
	WriteJSON(w, map[string]interface{}{
		"message": "Teams retrieved successfully",
		"teams":   teams,
	}, http.StatusOK)
}

func (h *Handler) GetTeamById(w http.ResponseWriter, r *http.Request) {
	idS := chi.URLParam(r, "team-id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		WriteJSON(w, map[string]interface{}{
			"message": "Could not convert id",
			"error":   err.Error(),
		}, http.StatusBadRequest)
		return
	}
	team, err := h.RepositoryService.ReadTeamByID(int64(id))
	if err != nil {
		WriteJSON(w, map[string]interface{}{
			"message": "Could not get the team",
			"error":   err.Error(),
		}, http.StatusNotFound)
		return
	}
	WriteJSON(w, map[string]interface{}{
		"message": "Team retrieved successfully",
		"team":    team,
	}, http.StatusOK)
}

func (h *Handler) ModifyTeam(w http.ResponseWriter, r *http.Request) {
	team := new(model.Team)
	err := json.NewDecoder(r.Body).Decode(&team)
	if err != nil {
		WriteJSON(w, map[string]interface{}{
			"message": "Invalid team data or bad format",
			"error":   err.Error(),
		}, http.StatusBadRequest)
		return
	}
	err = h.RepositoryService.UpdateTeam(team)
	if err != nil {
		WriteJSON(w, map[string]interface{}{
			"message": "Could not update team",
			"error":   err.Error(),
		}, http.StatusBadRequest)
		return
	}
	WriteJSON(w, map[string]interface{}{
		"message": "Team updated successfully",
		"team":    team,
	}, http.StatusOK)
}

func (h *Handler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	idS := chi.URLParam(r, "team-id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		WriteJSON(w, map[string]interface{}{
			"message": "Could not convert id",
			"error":   err.Error(),
		}, http.StatusBadRequest)
		return
	}
	err = h.RepositoryService.DeleteTeamByID(int64(id))
	if err != nil {
		WriteJSON(w, map[string]interface{}{
			"message": "Could not delete team",
			"error":   err.Error(),
		}, http.StatusBadRequest)
		return
	}
	WriteJSON(w, map[string]interface{}{
		"message": "Team deleted successfully",
	}, http.StatusOK)
}
