package main

import (
	"fmt"
	"net/http"

	"medic-api/clients"
	"medic-api/config"
	"medic-api/controllers"
	"medic-api/repositories"
	"medic-api/routes"
	"medic-api/services"
)

func main() {
	// =========================
	// SQL Server
	// =========================
	db, err := config.ConnectDB()
	if err != nil {
		panic(fmt.Sprintf("Database connection failed: %v", err))
	}

	defer db.Close()

	fmt.Println("SQL Server connected successfully")

	// =========================
	// Oracle
	// =========================
	oracleDB, err := config.ConnectOracle()
	if err != nil {
		panic(fmt.Sprintf("Oracle connection failed: %v", err))
	}

	defer oracleDB.Close()

	fmt.Println("Oracle connected successfully")

	// =========================
	// Repository
	// =========================
	healthRepository := repositories.NewHealthRepository(db)
	patientRepository := repositories.NewPatientRepository(db)

	// Room menggunakan Oracle
	roomRepository := repositories.NewRoomRepository(oracleDB)

	// =========================
	// Client
	// =========================
	aplicareClient := clients.NewAplicareClient()

	// =========================
	// Service
	// =========================
	healthService := services.NewHealthService(healthRepository)
	patientService := services.NewPatientService(patientRepository)

	// Room menggunakan Oracle + BPJS Aplicare
	roomService := services.NewRoomService(
		roomRepository,
		aplicareClient,
	)

	// =========================
	// Controller
	// =========================
	healthController := controllers.NewHealthController(healthService)
	patientController := controllers.NewPatientController(patientService)

	// Room menggunakan Oracle + BPJS Aplicare
	roomController := controllers.NewRoomController(roomService)

	// =========================
	// Router
	// =========================
	mux := http.NewServeMux()

	routes.RegisterRoutes(
		mux,
		healthController,
		patientController,
		roomController,
	)

	fmt.Println("MEDIC API running on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
