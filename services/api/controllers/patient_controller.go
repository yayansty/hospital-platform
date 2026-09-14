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
		)
		return
	}

	helpers.Success(
		w,
		http.StatusOK,
		"Patients retrieved successfully",
		patients,
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
		)
		return
	}

	if errors.Is(err, helpers.ErrDatabase) {
		helpers.Error(
			w,
			http.StatusInternalServerError,
			"Database error",
			nil,
		)
		return
	}

	helpers.Error(
		w,
		http.StatusInternalServerError,
		"Failed to get patient",
		nil,
	)
	return
}

	helpers.Success(
		w,
		http.StatusOK,
		"Patient retrieved successfully",
		patient,
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
		)
		return
	}

	err := c.service.CreatePatient(&patient)

	if err != nil {
		var validationErr *helpers.ValidationError

		if errors.As(err, &validationErr) {
			helpers.Error(
				w,
				http.StatusBadRequest,
				"Validation failed",
				validationErr.Errors,
			)
			return
		}

		helpers.Error(
			w,
			http.StatusInternalServerError,
			"Failed to create patient",
			nil,
		)
		return
	}

	helpers.Success(
		w,
		http.StatusCreated,
		"Patient created successfully",
		patient,
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
			)
			return
		}

		if errors.Is(err, helpers.ErrNotFound) {
	         helpers.Error(
		     w,
		     http.StatusNotFound,
		     "Patient not found",
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
	       )
	       return
        }

		helpers.Error(
			w,
			http.StatusInternalServerError,
			"Failed to update patient",
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
		)
		return
	}

	helpers.Success(
		w,
		http.StatusOK,
		"Patient updated successfully",
		updatedPatient,
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
		 )
		return
	}

	if errors.Is(err, helpers.ErrDatabase) {
		  helpers.Error(
			w,
			http.StatusInternalServerError,
			"Database error",
			nil,
		)
		return
	}

	helpers.Error(
		w,
		http.StatusInternalServerError,
		"Failed to delete patient",
		nil,
	)
	return
    }

	helpers.Success(
		w,
		http.StatusOK,
		"Patient deleted successfully",
		nil,
	)
}
