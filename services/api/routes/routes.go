package routes

import (
	"net/http"

	"medic-api/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	healthController *controllers.HealthController,
	patientController *controllers.PatientController,
	roomController *controllers.RoomController,
) {
	mux.HandleFunc("/api/health", healthController.Check)

	mux.HandleFunc("GET /api/patients", patientController.GetAll)
	mux.HandleFunc("GET /api/patients/{id}", patientController.GetByID)
	mux.HandleFunc("POST /api/patients", patientController.Create)
	mux.HandleFunc("PUT /api/patients/{id}", patientController.Update)
	mux.HandleFunc("DELETE /api/patients/{id}", patientController.Delete)

	mux.HandleFunc("GET /api/rooms/availability", roomController.GetAvailability)

	mux.HandleFunc("POST /api/rooms/availability/update-bpjs/{koderuang}", roomController.UpdateOneRoomToBPJS)
	mux.HandleFunc("GET /api/rooms/availability/bpjs", roomController.ReadRoomsFromBPJS)
	mux.HandleFunc("POST /api/rooms/availability/sync-bpjs", roomController.SyncRoomsFromBPJS)
}
