package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"medic-api/helpers"
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
		helpers.Error(
			w,
			http.StatusInternalServerError,
			"Failed to get patients",
			nil,
			nil,
		)
		return
	}

	helpers.Success(
		w,
		http.StatusOK,
		"Patients retrieved successfully",
		patients,
		nil,
	)
}

func (c *PatientController) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helpers.Error(
			w,
			http.StatusBadRequest,
			"Invalid patient ID",
			nil,
			nil,
		)
		return
	}

	patient, err := c.service.GetPatientByID(id)
	if err != nil {
		if errors.Is(err, helpers.ErrNotFound) {
			helpers.Error(
				w,
				http.StatusNotFound,
				"Patient not found",
				nil,
				nil,
			)
			return
		}

		if errors.Is(err, helpers.ErrDatabase) {
			helpers.Error(
				w,
				http.StatusInternalServerError,
				"Database error",
				nil,
				nil,
			)
			return
		}

		helpers.Error(
			w,
			http.StatusInternalServerError,
			"Failed to get patient",
			nil,
			nil,
		)
		return
	}

	helpers.Success(
		w,
		http.StatusOK,
		"Patient retrieved successfully",
		patient,
		nil,
	)
}

func (c *PatientController) Create(w http.ResponseWriter, r *http.Request) {
	var patient models.Patient

	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		helpers.Error(
			w,
			http.StatusBadRequest,
			"Invalid JSON",
			nil,
			nil,
		)
		return
	}

	//err := c.service.CreatePatient(&patient)

	// if err != nil {
	// 	var validationErr *helpers.ValidationError

	// 	if errors.As(err, &validationErr) {
	// 		helpers.Error(
	// 			w,
	// 			http.StatusBadRequest,
	// 			"Validation failed",
	// 			validationErr.Errors,
	// 		)
	// 		return
	// 	}

	// 	helpers.Error(
	// 		w,
	// 		http.StatusInternalServerError,
	// 		"Failed to create patient",
	// 		nil,
	// 	)
	// 	return
	// }
	if err := c.service.CreatePatient(&patient); err != nil {
		var validationErr *helpers.ValidationError

		if errors.As(err, &validationErr) {
			helpers.Error(
				w,
				http.StatusBadRequest,
				"Validation failed",
				validationErr.Errors,
				nil,
			)
			return
		}

		if errors.Is(err, helpers.ErrDuplicateMedicalRecordNumber) {
			helpers.Error(
				w,
				http.StatusConflict,
				"Medical record number already exists",
				nil,
				nil,
			)
			return
		}

		http.Error(w, "Failed to create patient", http.StatusInternalServerError)
		return
	}

	helpers.Success(
		w,
		http.StatusOK,
		"Patient created successfully",
		patient,
		nil,
	)
}

func (c *PatientController) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helpers.Error(
			w,
			http.StatusBadRequest,
			"Invalid patient ID",
			nil,
			nil,
		)
		return
	}

	var patient models.Patient

	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		helpers.Error(
			w,
			http.StatusBadRequest,
			"Invalid JSON",
			nil,
			nil,
		)
		return
	}

	patient.ID = id

	err = c.service.UpdatePatient(&patient)

	if err != nil {
		var validationErr *helpers.ValidationError

		if errors.As(err, &validationErr) {
			helpers.Error(
				w,
				http.StatusBadRequest,
				"Validation failed",
				validationErr.Errors,
				nil,
			)
			return
		}

		if errors.Is(err, helpers.ErrNotFound) {
			helpers.Error(
				w,
				http.StatusNotFound,
				"Patient not found",
				nil,
				nil,
			)
			return
		}

		if errors.Is(err, helpers.ErrDatabase) {
			helpers.Error(
				w,
				http.StatusInternalServerError,
				"Database error",
				nil,
				nil,
			)
			return
		}

		helpers.Error(
			w,
			http.StatusInternalServerError,
			"Failed to update patient",
			nil,
			nil,
		)
		return
	}

	updatedPatient, err := c.service.GetPatientByID(id)
	if err != nil {
		helpers.Error(
			w,
			http.StatusInternalServerError,
			"Failed to get updated patient",
			nil,
			nil,
		)
		return
	}

	helpers.Success(
		w,
		http.StatusOK,
		"Patient updated successfully",
		updatedPatient,
		nil,
	)
}

func (c *PatientController) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helpers.Error(
			w,
			http.StatusBadRequest,
			"Invalid patient ID",
			nil,
			nil,
		)
		return
	}

	if err := c.service.DeletePatient(id); err != nil {
		if errors.Is(err, helpers.ErrNotFound) {
			helpers.Error(
				w,
				http.StatusNotFound,
				"Patient not found",
				nil,
				nil,
			)
			return
		}

		if errors.Is(err, helpers.ErrDatabase) {
			helpers.Error(
				w,
				http.StatusInternalServerError,
				"Database error",
				nil,
				nil,
			)
			return
		}

		helpers.Error(
			w,
			http.StatusInternalServerError,
			"Failed to delete patient",
			nil,
			nil,
		)
		return
	}

	helpers.Success(
		w,
		http.StatusOK,
		"Patient deleted successfully",
		nil,
		nil,
	)
}

func (c *PatientController) Patch(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helpers.Error(
			w,
			http.StatusBadRequest,
			"Invalid patient ID",
			nil,
			nil,
		)
		return
	}

	var patch models.PatientPatch

	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		helpers.Error(
			w,
			http.StatusBadRequest,
			"Invalid JSON",
			nil,
			nil,
		)
		return
	}

	err = c.service.PatchPatient(id, &patch)

	if err != nil {
		var validationErr *helpers.ValidationError

		if errors.As(err, &validationErr) {
			helpers.Error(
				w,
				http.StatusBadRequest,
				"Validation failed",
				validationErr.Errors,
				nil,
			)
			return
		}

		if errors.Is(err, helpers.ErrNotFound) {
			helpers.Error(
				w,
				http.StatusNotFound,
				"Patient not found",
				nil,
				nil,
			)
			return
		}

		helpers.Error(
			w,
			http.StatusInternalServerError,
			"Failed to update patient",
			nil,
			nil,
		)
		return
	}

	patient, err := c.service.GetPatientByID(id)
	if err != nil {
		helpers.Error(
			w,
			http.StatusInternalServerError,
			"Failed to get updated patient",
			nil,
			nil,
		)
		return
	}

	helpers.Success(
		w,
		http.StatusOK,
		"Patient updated successfully",
		patient,
		nil,
	)
}
