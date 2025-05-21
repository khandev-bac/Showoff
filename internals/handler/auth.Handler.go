package handler

import (
	"encoding/json"
	"errors"
	"exceapp/internals/model"
	"exceapp/internals/service"
	"exceapp/pkg/google"
	"exceapp/pkg/jwt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Check(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ok"))
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	createdUser, err := h.service.Signup(ctx, user)
	if err != nil {
		http.Error(w, "Failed to create account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tokens, err := jwt.GenerateJWTToken(createdUser.ID)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	if err := h.service.UpdateRefreshToken(ctx, createdUser.ID, tokens.RefreshToken); err != nil {
		http.Error(w, "Could not save refresh token", http.StatusInternalServerError)
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
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
	})

	response := map[string]interface{}{
		"user_id":       createdUser.ID,
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Email    string `json:"user_email"`
		Password string `json:"user_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid login request", http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	tokens, err := jwt.GenerateJWTToken(user.ID)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	if err := h.service.UpdateRefreshToken(ctx, user.ID, tokens.RefreshToken); err != nil {
		http.Error(w, "Failed to update refresh token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(15 * time.Minute),
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
	})

	response := map[string]interface{}{
		"message": "Login successful",
		"user_id": user.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := uuid.NewString()
	url := google.GetLoginUrl(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
func (h *UserHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code", http.StatusBadRequest)
		return
	}
	token, err := google.ExchangeToken(code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}
	userInfo, err := google.GetUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	user, err := h.service.FindByEmail(ctx, userInfo.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		user, err = h.service.Signup(ctx, &model.User{
			Name:       userInfo.Name,
			Email:      userInfo.Email,
			ProfilePic: userInfo.ProfilePic,
		})
		if err != nil {
			http.Error(w, "User creation failed", http.StatusInternalServerError)
			return
		}
	}
	tokens, err := jwt.GenerateJWTToken(user.ID)
	if err != nil {
		http.Error(w, "JWT creation failed", http.StatusInternalServerError)
		return
	}
	if err := h.service.UpdateRefreshToken(ctx, user.ID, tokens.RefreshToken); err != nil {
		http.Error(w, "Refresh token save failed", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(3 * 30 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Login successful",
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})

}
