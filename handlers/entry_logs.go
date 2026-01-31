package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"iot_proj/models"
)

func (h *MQTTHandler) GetEntryLogs(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

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

	logs, err := h.Service.GetEntryLogsForUser(user.ID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(logs)
}
