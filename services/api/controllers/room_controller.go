package controllers

import (
	"net/http"

	"medic-api/helpers"
	"medic-api/services"
)

type RoomController struct {
	service *services.RoomService
}

func NewRoomController(service *services.RoomService) *RoomController {
	return &RoomController{service: service}
}

func (c *RoomController) GetAvailability(w http.ResponseWriter, r *http.Request) {
	rooms, err := c.service.GetRoomAvailability()
	if err != nil {
		helpers.Error(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve room availability",
			nil,
			nil,
		)
		return
	}

	helpers.Success(
		w,
		http.StatusOK,
		"Room availability retrieved successfully",
		rooms,
		nil,
	)
}

func (c *RoomController) UpdateOneRoomToBPJS(w http.ResponseWriter, r *http.Request) {
	kodeRuang := r.PathValue("koderuang")

	if kodeRuang == "" {
		helpers.Error(
			w,
			http.StatusBadRequest,
			"Kode ruang is required",
			nil,
			nil,
		)
		return
	}

	// Update semua room
	if kodeRuang == "all=true" {
		response, err := c.service.UpdateAllRoomsToBPJS()
		if err != nil {
			helpers.Error(
				w,
				http.StatusInternalServerError,
				"Failed to update all rooms to BPJS",
				err.Error(),
				nil,
			)
			return
		}

		helpers.Success(
			w,
			http.StatusOK,
			"All rooms successfully updated to BPJS",
			response,
		    nil,
		)
		return
	}

	// Update satu room
	response, err := c.service.UpdateOneRoomToBPJS(kodeRuang)
	if err != nil {
		helpers.Error(
			w,
			http.StatusInternalServerError,
			"Failed to update room to BPJS",
			err.Error(),
			nil,
		)
		return
	}

	helpers.Success(
		w,
		http.StatusOK,
		"Room successfully updated to BPJS",
		response,
		nil,
	)
}

func (c *RoomController) ReadRoomsFromBPJS(w http.ResponseWriter, r *http.Request) {
	response, err := c.service.ReadRoomsFromBPJS(1, 20)
	if err != nil {
		helpers.Error(
			w,
			http.StatusInternalServerError,
			"Failed to read rooms from BPJS",
			err.Error(),
			nil,
		)
		return
	}

	helpers.Success(
		w,
		http.StatusOK,
		"Rooms successfully retrieved from BPJS",
		response,
		nil,
	)
}

func (c *RoomController) SyncRoomsFromBPJS(
	w http.ResponseWriter,
	r *http.Request,
) {
	results, err := c.service.SyncRoomsFromBPJS()

	if err != nil {
		helpers.Error(
			w,
			http.StatusInternalServerError,
			"Failed to sync rooms from BPJS",
			err.Error(),
			nil,
		)
		return
	}

	helpers.Success(
		w,
		http.StatusOK,
		"Rooms successfully synchronized from BPJS to Oracle",
		results,
		nil,
	)
}
