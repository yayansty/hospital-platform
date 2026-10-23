package routes

import (
	"net/http"

	"medic-api/controllers"
	"medic-api/middleware"
)

func RegisterRoutes(
	mux *http.ServeMux,
	healthController *controllers.HealthController,
	patientController *controllers.PatientController,
	roomController *controllers.RoomController,
	authController *controllers.AuthController,
	profileController *controllers.ProfileController,
) {
	mux.HandleFunc("/api/health", healthController.Check)

	mux.HandleFunc("POST /api/auth/login", authController.Login)
	mux.HandleFunc("GET /api/patients", patientController.GetAll)
	mux.HandleFunc("GET /api/patients/{id}", patientController.GetByID)
	mux.HandleFunc("POST /api/patients", patientController.Create)
	mux.HandleFunc("PUT /api/patients/{id}", patientController.Update)
	mux.HandleFunc("PATCH /api/patients/{id}", patientController.Patch)
	mux.HandleFunc("DELETE /api/patients/{id}", patientController.Delete)

	mux.HandleFunc("GET /api/rooms/availability", roomController.GetAvailability)

	mux.Handle("GET /api/auth/me", middleware.JWTAuth(http.HandlerFunc(profileController.Me)))

	mux.HandleFunc("POST /api/rooms/availability/update-bpjs/{koderuang}", roomController.UpdateOneRoomToBPJS)
	mux.HandleFunc("GET /api/rooms/availability/bpjs", roomController.ReadRoomsFromBPJS)
	mux.HandleFunc("POST /api/rooms/availability/sync-bpjs", roomController.SyncRoomsFromBPJS)
}
