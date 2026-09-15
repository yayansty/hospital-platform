package services

import (
	"fmt"

	"medic-api/clients"
	"medic-api/models"
	"medic-api/repositories"
)

type RoomService struct {
	repository     *repositories.RoomRepository
	aplicareClient *clients.AplicareClient
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

func (s *RoomService) UpdateAllRoomsToBPJS() ([]*clients.AplicareResponse, error) {
	rooms, err := s.repository.FindAvailability()
	if err != nil {
		return nil, err
	}

	var responses []*clients.AplicareResponse

	for _, room := range rooms {
		response, err := s.aplicareClient.UpdateBed(room)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to update room %s to BPJS: %w",
				room.KodeRuang,
				err,
			)
		}

		responses = append(responses, response)
	}

	return responses, nil
}
