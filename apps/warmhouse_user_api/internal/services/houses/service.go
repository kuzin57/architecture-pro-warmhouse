package houses

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/warmhouse/warmhouse_user_api/internal/entities"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/async"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"

	"github.com/warmhouse/libraries/convert"

	"github.com/google/uuid"
)

type Service struct {
	housesRepository  HousesRepository
	devicesRepository DevicesRepository
	broker            *async.UserController
}

func NewService(housesRepository HousesRepository, devicesRepository DevicesRepository, broker *async.UserController) *Service {
	return &Service{housesRepository: housesRepository, devicesRepository: devicesRepository, broker: broker}
}

func (s *Service) GetUserHouses(ctx context.Context, userID uuid.UUID) ([]server.House, error) {
	houses, err := s.housesRepository.GetUserHouses(ctx, userID)
	if err != nil {
		return nil, err
	}

	return convert.MapSlice(houses, func(house entities.House) server.House {
		return server.House{
			Id:        house.ID,
			Name:      house.Name,
			Address:   house.Address,
			CreatedAt: house.CreatedAt,
			UpdatedAt: house.UpdatedAt,
		}
	}), nil
}

func (s *Service) CreateHouse(ctx context.Context, request server.CreateHouseRequestObject) (uuid.UUID, error) {
	house := entities.House{
		ID:      uuid.New(),
		UserID:  request.Params.XUserId,
		Name:    request.Body.Name,
		Address: request.Body.Address,
	}

	err := s.housesRepository.CreateHouse(ctx, house)
	if err != nil {
		return uuid.UUID{}, err
	}

	return house.ID, nil
}

func (s *Service) GetHouseInfo(ctx context.Context, houseID, userID uuid.UUID) (server.House, error) {
	house, err := s.housesRepository.GetHouse(ctx, houseID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server.House{}, ErrHouseNotFound
		}
		return server.House{}, err
	}

	return server.House{
		Id:        house.ID,
		Name:      house.Name,
		Address:   house.Address,
		CreatedAt: house.CreatedAt,
		UpdatedAt: house.UpdatedAt,
	}, nil
}

func (s *Service) UpdateHouse(ctx context.Context, request server.UpdateHouseRequestObject) (server.House, error) {
	house, err := s.housesRepository.GetHouse(ctx, request.Body.HouseId, request.Params.XUserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server.House{}, ErrHouseNotFound
		}
		return server.House{}, err
	}

	house.Name = convert.UnwrapOr(request.Body.Name, house.Name)
	house.Address = convert.UnwrapOr(request.Body.Address, house.Address)

	err = s.housesRepository.UpdateHouse(ctx, house)
	if err != nil {
		return server.House{}, err
	}

	return server.House{
		Id:        house.ID,
		Name:      house.Name,
		Address:   house.Address,
		CreatedAt: house.CreatedAt,
		UpdatedAt: house.UpdatedAt,
	}, nil
}

func (s *Service) DeleteHouse(ctx context.Context, houseID, userID uuid.UUID) error {
	_, err := s.housesRepository.GetHouse(ctx, houseID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrHouseNotFound
		}
		return err
	}

	devices, err := s.devicesRepository.GetHouseDevices(ctx, houseID, userID)
	if err != nil {
		return fmt.Errorf("error getting house devices: %w", err)
	}

	for _, device := range devices {
		deviceDeletedMessage := async.NewDeviceDeletedMessageFromDeviceDeletedChannel()
		deviceDeletedMessage.Payload.DeviceId = device.ID.String()
		deviceDeletedMessage.Payload.HouseId = houseID.String()

		if err := s.broker.SendToDeviceDeletedOperation(ctx, deviceDeletedMessage); err != nil {
			return fmt.Errorf("error sending device deleted message: %w", err)
		}
	}

	err = s.housesRepository.DeleteHouse(ctx, houseID)
	if err != nil {
		return fmt.Errorf("error deleting house: %w", err)
	}

	return nil
}
