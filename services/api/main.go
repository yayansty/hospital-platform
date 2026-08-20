package main

import (
	"fmt"
	"net/http"

	"medic-api/config"
	"medic-api/controllers"
	"medic-api/repositories"
	"medic-api/routes"
	"medic-api/services"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		panic(fmt.Sprintf("Database connection failed: %v", err))
	}

	defer db.Close()

	fmt.Println("SQL Server connected successfully")

	// Repository
	healthRepository := repositories.NewHealthRepository(db)
	patientRepository := repositories.NewPatientRepository(db)

	// Service
	healthService := services.NewHealthService(healthRepository)
	patientService := services.NewPatientService(patientRepository)

	// Controller
	healthController := controllers.NewHealthController(healthService)
	patientController := controllers.NewPatientController(patientService)

	// Router
	mux := http.NewServeMux()

	routes.RegisterRoutes(
		mux,
		healthController,
		patientController,
	)

	fmt.Println("MEDIC API running on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
