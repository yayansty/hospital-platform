package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"medic-api/models"
	"medic-api/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format request tidak valid", http.StatusBadRequest)
		return
	}

	result, err := c.authService.Login(
		r.Context(),
		req.Username,
		req.Password,
	)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			http.Error(w, "Username atau password salah", http.StatusUnauthorized)
		case errors.Is(err, services.ErrInactiveUser):
			http.Error(w, "Akun pengguna tidak aktif", http.StatusUnauthorized)
		case errors.Is(err, services.ErrJWTSecretMissing):
			http.Error(w, "Konfigurasi autentikasi bermasalah", http.StatusInternalServerError)
		default:
			http.Error(w, "Terjadi kesalahan pada server", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Login berhasil",
		"data":    result,
	})
}
