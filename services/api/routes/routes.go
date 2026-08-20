package routes

import (
	"net/http"

	"medic-api/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	healthController *controllers.HealthController,
	patientController *controllers.PatientController,
) {
	mux.HandleFunc("/api/health", healthController.Check)

	mux.HandleFunc("GET /api/patients", patientController.GetAll)
	mux.HandleFunc("GET /api/patients/{id}", patientController.GetByID)
	mux.HandleFunc("POST /api/patients", patientController.Create)
	mux.HandleFunc("PUT /api/patients/{id}", patientController.Update)
	mux.HandleFunc("DELETE /api/patients/{id}", patientController.Delete)
}
