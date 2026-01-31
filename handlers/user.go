package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"iot_proj/models"
)

func (h *MQTTHandler) UpdateUserKeyCard(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		KeyCard string `json:"key_card"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateUserKeyCard(userID, requestBody.KeyCard); err != nil {
		http.Error(w, "Failed to update key card", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Key card updated successfully"})
}

func (h *MQTTHandler) UpdateUserAccessLevel(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		AccessLevel int `json:"access_level"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateUserAccessLevel(userID, requestBody.AccessLevel); err != nil {
		http.Error(w, "Failed to update access level", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Access level updated successfully"})
}

func (h *MQTTHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	users, err := h.Service.GetUsers(limit, offset)
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (h *MQTTHandler) Profile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}
