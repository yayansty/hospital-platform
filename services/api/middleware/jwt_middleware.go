package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const authClaimsKey contextKey = "authClaims"

type AuthClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv("JWT_SECRET")
		if len(secret) == 0 {
			writeAuthError(w, http.StatusInternalServerError, "Konfigurasi autentikasi bermasalah")
			return
		}

		authHeader := r.Header.Get("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") ||
			strings.TrimSpace(parts[1]) == "" {
			writeAuthError(w, http.StatusUnauthorized, "Token autentikasi diperlukan")
			return
		}

		claims := &AuthClaims{}

		token, err := jwt.ParseWithClaims(
			parts[1],
			claims,
			func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			},
			jwt.WithIssuer("medic-api"),
		)

		if err != nil || token == nil || !token.Valid {
			writeAuthError(w, http.StatusUnauthorized, "Token tidak valid atau kedaluwarsa")
			return
		}

		ctx := context.WithValue(r.Context(), authClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetAuthClaims(r *http.Request) (*AuthClaims, bool) {
	claims, ok := r.Context().Value(authClaimsKey).(*AuthClaims)
	return claims, ok
}

func RequireRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetAuthClaims(r)
			if !ok {
				writeAuthError(w, http.StatusUnauthorized, "Autentikasi diperlukan")
				return
			}

			for _, role := range allowedRoles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			writeAuthError(w, http.StatusForbidden, "Anda tidak memiliki hak akses")
		}))
	}
}

func writeAuthError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    status,
		"message": message,
	})
}
