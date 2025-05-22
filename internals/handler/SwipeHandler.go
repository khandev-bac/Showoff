package handler

import (
	"encoding/json"
	"exceapp/internals/service"
	"net/http"

	"github.com/google/uuid"
)

type SwipeHandler struct {
	service *service.SwipeService
}

func NewSwipeHandler(service *service.SwipeService) *SwipeHandler {
	return &SwipeHandler{
		service: service,
	}
}

// /api/swipe/next
func (hs *SwipeHandler) GetNextUnswipedUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDStr := r.Context().Value("USERID").(string)
	if userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "user id is not valid", http.StatusUnauthorized)
		return
	}
	user, err := hs.service.GetUnswippedUsers(ctx, userID)
	if err != nil {
		http.Error(w, "No profiles left", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// POST: api/swipe/
func (hs *SwipeHandler) SwipeUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDStr := r.Context().Value("USERID").(string)
	if userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "user id is not valid", http.StatusUnauthorized)
		return
	}
	var req struct {
		SwipedID uuid.UUID `json:"swiped_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	if err := hs.service.SaveSwipes(ctx, userID, req.SwipedID); err != nil {
		http.Error(w, "Failed to swipe", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Swipe saved"})
}
func (hs *SwipeHandler) GetSwippedHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDStr := ctx.Value("USERID").(string)
	if userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "user id is not valid", http.StatusUnauthorized)
		return
	}
	users, err := hs.service.GetSwippedHistory(ctx, userID)
	if err != nil {
		http.Error(w, "Failed to fetch swipe history", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}
