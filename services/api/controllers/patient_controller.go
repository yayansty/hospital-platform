package controllers

import (
	"encoding/json"
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

// func (c *PatientController) GetAll(w http.ResponseWriter, r *http.Request) {
// 	patients, err := c.service.GetAllPatients()
// 	if err != nil {
// 		helpers.HandleError(w, err)
// 		return
// 	}

// 	helpers.Success(
// 		w,
// 		http.StatusOK,
// 		"Patients retrieved successfully",
// 		patients,
// 		nil,
// 	)
// }

func (c *PatientController) GetAll(
	w http.ResponseWriter,
	r *http.Request,
) {
	query := r.URL.Query()

	page := 0
	limit := 0
	var err error

	if value := query.Get("page"); value != "" {
		page, err = strconv.Atoi(value)
		if err != nil {
			helpers.Error(
				w, http.StatusBadRequest,
				"Invalid page parameter", nil, nil,
			)
			return
		}
	}

	if value := query.Get("limit"); value != "" {
		limit, err = strconv.Atoi(value)
		if err != nil {
			helpers.Error(
				w, http.StatusBadRequest,
				"Invalid limit parameter", nil, nil,
			)
			return
		}
	}

	// filter := models.PatientFilter{
	// 	Page:                page,
	// 	Limit:               limit,
	// 	Search:              query.Get("search"),
	// 	MedicalRecordNumber: query.Get("medical_record_number"),
	// 	Gender:              query.Get("gender"),
	// }

	filter := models.PatientFilter{
		Page:                page,
		Limit:               limit,
		Search:              query.Get("search"),
		MedicalRecordNumber: query.Get("medical_record_number"),
		Gender:              query.Get("gender"),
		SortBy:              query.Get("sort_by"),
		SortOrder:           query.Get("sort_order"),
		DateFrom:            query.Get("date_from"),
		DateTo:              query.Get("date_to"),
	}

	patients, meta, err := c.service.GetAllPatients(filter)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	helpers.Success(
		w,
		http.StatusOK,
		"Patients retrieved successfully",
		patients,
		meta,
	)
}

func (c *PatientController) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	patient, err := c.service.GetPatientByID(id)
	if err != nil {
		helpers.HandleError(w, err)
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
		helpers.HandleError(w, err)
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
		helpers.HandleError(w, err)
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
		helpers.HandleError(w, err)
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
		helpers.HandleError(w, err)
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
