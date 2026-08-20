package controllers

import (
	"encoding/json"
	"net/http"

	"medic-api/services"
)

type HealthController struct {
	Service *services.HealthService
}

func NewHealthController(service *services.HealthService) *HealthController {
	return &HealthController{
		Service: service,
	}
}

func (c *HealthController) Check(w http.ResponseWriter, r *http.Request) {
	err := c.Service.CheckDatabase()

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Database connection failed",
		})

		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"message":  "MEDIC API is running",
		"database": "connected",
	})
}
