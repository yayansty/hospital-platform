package controllers

import (
	"encoding/json"
	"net/http"

	"medic-api/middleware"
)

type ProfileController struct{}

func NewProfileController() *ProfileController {
	return &ProfileController{}
}

func (c *ProfileController) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetAuthClaims(r)
	if !ok {
		http.Error(w, "Autentikasi diperlukan", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Profil berhasil diambil",
		"data": map[string]interface{}{
			"user_id":  claims.UserID,
			"username": claims.Username,
			"role":     claims.Role,
		},
	})
}
