package handler

import (
	"context"
	"encoding/json"
	"exceapp/internals/model"
	"exceapp/internals/service"
	"exceapp/pkg/google"
	"exceapp/pkg/jwt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var ctx = context.Background()

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}
func (h *UserHandler) Check(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"user_name"`
		Email    string `json:"user_email"`
		Password string `json:"user_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Unable to take your request", http.StatusBadRequest)
		return
	}
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	tokens, err := jwt.GenerateJWTToken(user.ID)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	if err := h.service.UpdateRefreshToken(ctx, user.ID, tokens.RefreshToken); err != nil {
		http.Error(w, "Failed to update token", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(15 * time.Minute),
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(3 * 30 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
	})
	response := map[string]interface{}{
		"user_id":       user.ID,
		"user_name":     user.Name,
		"user_email":    user.Email,
		"token":         tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"user_email"`
		Password string `json:"user_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to process request", http.StatusBadRequest)
		return
	}
	user, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusBadRequest)
		return
	}
	tokens, err := jwt.GenerateJWTToken(user.ID)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	if err := h.service.UpdateRefreshToken(ctx, user.ID, tokens.RefreshToken); err != nil {
		http.Error(w, "Failed to update token", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(15 * time.Minute),
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(3 * 30 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": user.ID,
		"message": "Login successfull 🎉",
	})
}

func (h *UserHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()
	url := google.GetLoginUrl(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
func (h *UserHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code in callback", http.StatusBadRequest)
		return
	}
	token, err := google.ExchangeToken(code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}
	googleUser, err := google.GetUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		return
	}
	user, err := h.service.FindByEmail(r.Context(), googleUser.Email)
	if err != nil || user == nil {
		user = &model.User{
			ID:           uuid.New(),
			Name:         googleUser.Name,
			Email:        googleUser.Email,
			ProfilePic:   googleUser.ProfilePic,
			IsOauthUser:  true,
			RefreshToken: token.RefreshToken,
			GoogleID:     googleUser.ID,
			CreatedAt:    time.Now(),
		}
		err = h.service.Signup(r.Context(), user)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	} else {
		if token.RefreshToken != "" {
			_ = h.service.UpdateRefreshToken(r.Context(), user.ID, token.RefreshToken)
		}
	}
	tokens, err := jwt.GenerateJWTToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Login successful",
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}
