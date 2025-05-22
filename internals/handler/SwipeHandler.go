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
