package services

import (
	"fmt"
	"time"

	"medic-api/clients"
	"medic-api/models"
	"medic-api/repositories"
)

type RoomService struct {
	repository     *repositories.RoomRepository
	aplicareClient *clients.AplicareClient
}

type BPJSSyncResult struct {
	KodeRuang string `json:"koderuang"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}

func NewRoomService(
	repository *repositories.RoomRepository,
	aplicareClient *clients.AplicareClient,
) *RoomService {
	return &RoomService{
		repository:     repository,
		aplicareClient: aplicareClient,
	}
}

func (s *RoomService) GetRoomAvailability() ([]models.RoomAvailability, error) {
	return s.repository.FindAvailability()
}

func (s *RoomService) UpdateOneRoomToBPJS(
	kodeRuang string,
) (*clients.AplicareResponse, error) {

	rooms, err := s.repository.FindAvailability()
	if err != nil {
		return nil, err
	}

	for _, room := range rooms {
		if room.KodeRuang == kodeRuang {
			return s.aplicareClient.UpdateBed(room)
		}
	}

	return nil, fmt.Errorf("room %s not found", kodeRuang)
}

// func (s *RoomService) UpdateAllRoomsToBPJS() ([]*clients.AplicareResponse, error) {
// 	rooms, err := s.repository.FindAvailability()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var responses []*clients.AplicareResponse

// 	for _, room := range rooms {
// 		response, err := s.aplicareClient.UpdateBed(room)
// 		if err != nil {
// 			return nil, fmt.Errorf(
// 				"failed to update room %s to BPJS: %w",
// 				room.KodeRuang,
// 				err,
// 			)
// 		}

// 		responses = append(responses, response)
// 	}

//		return responses, nil
//	}
func (s *RoomService) UpdateAllRoomsToBPJS() ([]*clients.AplicareResponse, error) {

	rooms, err := s.repository.FindAvailability()
	if err != nil {
		return nil, err
	}

	responses := make([]*clients.AplicareResponse, 0, len(rooms))

	for i, room := range rooms {

		fmt.Printf(
			"Updating room %d/%d to BPJS: %s\n",
			i+1,
			len(rooms),
			room.KodeRuang,
		)

		response, err := s.aplicareClient.UpdateBed(room)

		if err != nil {
			return nil, fmt.Errorf(
				"failed updating room %s: %w",
				room.KodeRuang,
				err,
			)
		}

		responses = append(responses, response)

		// =====================================================
		// Jeda antar request ke BPJS.
		// Jangan melakukan request 11 kamar secara berturut-turut
		// tanpa jeda.
		// =====================================================
		if i < len(rooms)-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return responses, nil
}

func (s *RoomService) ReadRoomsFromBPJS(
	start, limit int,
) (*clients.AplicareBedReadResponse, error) {
	return s.aplicareClient.ReadBeds(start, limit)
}

func (s *RoomService) SyncRoomsFromBPJS() ([]BPJSSyncResult, error) {

	response, err := s.aplicareClient.ReadBeds(1, 100)
	if err != nil {
		return nil, err
	}

	if response.MetaData.Code != 1 {
		return nil, fmt.Errorf(
			"BPJS error: %s",
			response.MetaData.Message,
		)
	}

	results := make([]BPJSSyncResult, 0, len(response.Response.List))

	for _, bed := range response.Response.List {

		err := s.repository.UpdateBPJSKamar(bed)

		if err != nil {
			results = append(results, BPJSSyncResult{
				KodeRuang: bed.KodeRuang,
				Success:   false,
				Message:   err.Error(),
			})
			continue
		}

		results = append(results, BPJSSyncResult{
			KodeRuang: bed.KodeRuang,
			Success:   true,
			Message:   "Berhasil update Oracle",
		})
	}

	return results, nil
}
