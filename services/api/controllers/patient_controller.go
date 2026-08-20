package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"medic-api/models"
	"medic-api/services"
)

type PatientController struct {
	service *services.PatientService
}

func NewPatientController(service *services.PatientService) *PatientController {
	return &PatientController{
		service: service,
	}
}

func (c *PatientController) GetAll(w http.ResponseWriter, r *http.Request) {
	patients, err := c.service.GetAllPatients()
	if err != nil {
		http.Error(w, "Failed to get patients", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    patients,
	})
}

func (c *PatientController) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	patient, err := c.service.GetPatientByID(id)
	if err != nil {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    patient,
	})
}

func (c *PatientController) Create(w http.ResponseWriter, r *http.Request) {
	var patient models.Patient

	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := c.service.CreatePatient(&patient); err != nil {
		http.Error(w, "Failed to create patient", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Patient created successfully",
		"data":    patient,
	})
}
func (c *PatientController) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	var patient models.Patient

	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	patient.ID = id

	if err := c.service.UpdatePatient(&patient); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Patient not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to update patient", http.StatusInternalServerError)
		return
	}

	updatedPatient, err := c.service.GetPatientByID(id)
	if err != nil {
		http.Error(w, "Failed to get updated patient", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Patient updated successfully",
		"data":    updatedPatient,
	})
}
func (c *PatientController) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	if err := c.service.DeletePatient(id); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Patient not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to delete patient", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Patient deleted successfully",
	})
}
