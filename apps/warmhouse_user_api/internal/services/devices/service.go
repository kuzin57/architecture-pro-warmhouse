package devices

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/warmhouse/warmhouse_user_api/internal/entities"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/async"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"

	"github.com/warmhouse/libraries/convert"

	"github.com/google/uuid"
)

type Service struct {
	devicesRepository   DevicesRepository
	houseRepository     HouseRepository
	telemetryRepository TelemetryRepository
	broker              *async.UserController
}

func NewService(
	devicesRepository DevicesRepository,
	houseRepository HouseRepository,
	telemetryRepository TelemetryRepository,
	broker *async.UserController,
) *Service {
	s := &Service{
		devicesRepository:   devicesRepository,
		houseRepository:     houseRepository,
		telemetryRepository: telemetryRepository,
		broker:              broker,
	}
	return s
}

func (s *Service) GetUserDevices(ctx context.Context, userID uuid.UUID) ([]server.Device, error) {
	devices, err := s.devicesRepository.GetUserDevices(ctx, userID)
	if err != nil {
		return nil, err
	}

	return convert.MapSlice(devices, func(device entities.Device) server.Device {
		return server.Device{
			Id:        device.ID,
			Name:      device.Name,
			Unit:      device.Unit,
			Value:     device.Value,
			Status:    string(device.Status),
			CreatedAt: device.CreatedAt,
			UpdatedAt: device.UpdatedAt,
		}
	}), nil
}

func (s *Service) GetHouseDevices(ctx context.Context, houseID, userID uuid.UUID) ([]server.Device, error) {
	devices, err := s.devicesRepository.GetHouseDevices(ctx, houseID, userID)
	if err != nil {
		return nil, err
	}

	return convert.MapSlice(devices, func(device entities.Device) server.Device {
		return server.Device{
			Id:        device.ID,
			Name:      device.Name,
			Unit:      device.Unit,
			Value:     device.Value,
			Status:    string(device.Status),
			CreatedAt: device.CreatedAt,
			UpdatedAt: device.UpdatedAt,
		}
	}), nil
}

func (s *Service) GetDeviceInfo(ctx context.Context, deviceID, userID uuid.UUID) (server.Device, error) {
	device, err := s.devicesRepository.GetDevice(ctx, deviceID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server.Device{}, ErrDeviceNotFound
		}

		return server.Device{}, err
	}

	return server.Device{
		Id:        device.ID,
		Name:      device.Name,
		Unit:      device.Unit,
		Value:     device.Value,
		Status:    string(device.Status),
		CreatedAt: device.CreatedAt,
		UpdatedAt: device.UpdatedAt,
	}, nil
}

func (s *Service) UpdateDevice(ctx context.Context, request server.UpdateDeviceRequestObject) (uuid.UUID, error) {
	device, err := s.devicesRepository.GetDevice(ctx, request.Body.DeviceId, request.Params.XUserId)
	if err != nil {
		return uuid.UUID{}, err
	}

	msg := async.NewDeviceUpdatedMessageFromDeviceUpdatedChannel()

	msg.Payload.DeviceId = request.Body.DeviceId.String()
	msg.Payload.HouseId = device.HouseID.String()
	msg.Payload.Name = request.Body.Name
	msg.Payload.Type = string(device.Type)
	msg.Payload.Unit = request.Body.Unit
	msg.Payload.Value = request.Body.Value
	msg.Payload.SensorId = convert.FromIntToInt64(request.Body.SensorId)
	msg.Payload.Location = request.Body.Location
	msg.Payload.Host = request.Body.Host
	msg.Payload.Schedule = request.Body.Schedule

	if err := s.broker.SendToDeviceUpdatedOperation(ctx, msg); err != nil {
		return uuid.UUID{}, err
	}

	return request.Body.DeviceId, nil
}

func (s *Service) DeleteDevice(ctx context.Context, deviceID, userID uuid.UUID) error {
	device, err := s.devicesRepository.GetDevice(ctx, deviceID, userID)
	if err != nil {
		return err
	}

	msg := async.NewDeviceDeletedMessageFromDeviceDeletedChannel()
	msg.Payload.DeviceId = deviceID.String()
	msg.Payload.HouseId = device.HouseID.String()
	msg.Payload.Type = string(device.Type)

	if err := s.broker.SendToDeviceDeletedOperation(ctx, msg); err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateDevice(ctx context.Context, request server.CreateDeviceRequestObject) (uuid.UUID, error) {
	_, err := s.houseRepository.GetHouse(ctx, request.Body.HouseId, request.Params.XUserId)
	if err != nil {
		return uuid.UUID{}, err
	}

	deviceID := uuid.New()

	deviceCreatedMessage := async.NewDeviceCreatedMessageFromDeviceCreatedChannel()
	deviceCreatedMessage.Payload.DeviceId = deviceID.String()
	deviceCreatedMessage.Payload.HouseId = request.Body.HouseId.String()
	deviceCreatedMessage.Payload.DeviceType = string(request.Body.Type)
	deviceCreatedMessage.Payload.Name = request.Body.Name
	deviceCreatedMessage.Payload.Unit = request.Body.Unit
	deviceCreatedMessage.Payload.Value = &request.Body.Value
	deviceCreatedMessage.Payload.Location = &request.Body.Location
	deviceCreatedMessage.Payload.SensorId = convert.FromIntToInt64(request.Body.SensorId)
	deviceCreatedMessage.Payload.Host = request.Body.Host
	deviceCreatedMessage.Payload.Schedule = request.Body.Schedule

	log.Println("deviceCreatedMessage", deviceCreatedMessage)

	if err := s.broker.SendToDeviceCreatedOperation(ctx, deviceCreatedMessage); err != nil {
		return uuid.UUID{}, err
	}

	return deviceID, nil
}

func (s *Service) GetDeviceTelemetry(ctx context.Context, deviceID, userID uuid.UUID, limit, offset int, fromDate, toDate *time.Time) ([]server.Telemetry, error) {
	_, err := s.devicesRepository.GetDevice(ctx, deviceID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDeviceNotFound
		}
		return nil, err
	}

	telemetry, err := s.telemetryRepository.GetDeviceTelemetry(ctx, deviceID, limit, offset, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	return convert.MapSlice(telemetry, func(t entities.Telemetry) server.Telemetry {
		return server.Telemetry{
			DeviceId:  t.DeviceID,
			Timestamp: t.Timestamp,
			Data:      t.Data,
		}
	}), nil
}
