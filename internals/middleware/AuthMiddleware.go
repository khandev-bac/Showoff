package middleware

import (
	"context"
	"exceapp/pkg/jwt"
	"net/http"
	"strings"
)

var JWT_KEY = []byte("JWT_KEY")

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil || strings.TrimSpace(cookie.Value) == "" {
			http.Error(w, "Unauthorized: No access token", http.StatusUnauthorized)
			return
		}
		claims, err := jwt.ValidateToken(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}
		userID, ok := claims["userID"].(string)
		if !ok || userID == "" {
			http.Error(w, "Unauthorized: Invalid claims", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "USERID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
