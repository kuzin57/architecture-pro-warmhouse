package gates

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/warmhouse/libraries/convert"
	"github.com/warmhouse/libraries/scheduler"
	gatesapi "github.com/warmhouse/warmhouse_devices/internal/clients/gates_api"
	"github.com/warmhouse/warmhouse_devices/internal/consts"
	"github.com/warmhouse/warmhouse_devices/internal/cron"
	"github.com/warmhouse/warmhouse_devices/internal/entities"
	"github.com/warmhouse/warmhouse_devices/internal/generated/async"
	"github.com/warmhouse/warmhouse_devices/internal/repositories/telemetry"
)

type Service struct {
	devicesRepository   DevicesRepository
	telemetryRepository *telemetry.Repository
	gatesClient         GatesClient
	cronScheduler       *scheduler.Scheduler
}

func NewService(
	devicesRepository DevicesRepository,
	telemetryRepository *telemetry.Repository,
	gatesClient GatesClient,
	cronScheduler *scheduler.Scheduler,
) *Service {
	return &Service{
		devicesRepository:   devicesRepository,
		telemetryRepository: telemetryRepository,
		gatesClient:         gatesClient,
		cronScheduler:       cronScheduler,
	}
}

func (s *Service) CreateDevice(ctx context.Context, payload async.DeviceCreatedMessageFromDeviceCreatedChannelPayload) error {
	var (
		device = entities.Device{
			ID:       uuid.MustParse(payload.DeviceId),
			HouseID:  uuid.MustParse(payload.HouseId),
			Name:     payload.Name,
			Host:     payload.Host,
			Unit:     consts.DeviceUnitGates,
			Status:   consts.DeviceStatusUnknown,
			Type:     consts.DeviceTypeGates,
			Value:    convert.UnwrapOr(payload.Value, "closed"),
			Location: convert.UnwrapOr(payload.Location, ""),
			SensorID: payload.SensorId,
			Schedule: convert.UnwrapOr(payload.Schedule, checkGatesStatusDefaultSchedule),
		}
		checkGatesStatusTask = cron.NewCheckGatesStatusTask(s.devicesRepository, s.gatesClient, s.telemetryRepository, device.ID)
	)

	if err := s.cronScheduler.AddJob(ctx, device.Schedule, device.ID, func(ctx context.Context) error {
		if err := checkGatesStatusTask.Run(ctx); err != nil {
			log.Println("error running check gates status task", err)

			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return s.devicesRepository.CreateDevice(ctx, device)
}

func (s *Service) UpdateDevice(ctx context.Context, payload async.DeviceUpdatedMessageFromDeviceUpdatedChannelPayload) error {
	deviceID, err := uuid.Parse(payload.DeviceId)
	if err != nil {
		return err
	}

	device, err := s.devicesRepository.GetDevice(ctx, deviceID)
	if err != nil {
		return err
	}

	device.Name = convert.UnwrapOr(payload.Name, device.Name)
	device.Unit = convert.UnwrapOr(payload.Unit, device.Unit)
	device.Value = convert.UnwrapOr(payload.Value, device.Value)
	device.Status = convert.UnwrapOr((*consts.DeviceStatus)(payload.Status), device.Status)
	device.Location = convert.UnwrapOr(payload.Location, device.Location)
	device.Host = convert.UnwrapOr(payload.Host, device.Host)
	device.SensorID = payload.SensorId
	device.Schedule = convert.UnwrapOr(payload.Schedule, device.Schedule)

	switch {
	case payload.Status != nil && *payload.Status != string(device.Status) && *payload.Status == string(consts.DeviceStatusActive):
		err = s.gatesClient.ActivateGate(ctx, device.Host)
		if err != nil {
			return err
		}

		if err := s.addGatesStatusTask(ctx, device); err != nil {
			return err
		}
	case payload.Status != nil && *payload.Status != string(device.Status) && *payload.Status == string(consts.DeviceStatusInactive):
		if err := s.cronScheduler.RemoveJob(ctx, device.ID); err != nil {
			return err
		}

		err = s.gatesClient.DeactivateGate(ctx, device.Host)
		if err != nil {
			return err
		}
	case payload.Value != nil && *payload.Value != string(device.Value) && *payload.Value == consts.GatesValueOpened:
		err = s.gatesClient.ChangeGateState(ctx, device.Host, gatesapi.GateStateOpen)
		if err != nil {
			return err
		}
	case payload.Value != nil && *payload.Value != string(device.Value) && *payload.Value == consts.GatesValueClosed:
		err = s.gatesClient.ChangeGateState(ctx, device.Host, gatesapi.GateStateClosed)
		if err != nil {
			return err
		}
	}

	if payload.Schedule != nil && *payload.Schedule != device.Schedule {
		if err := s.cronScheduler.RemoveJob(ctx, device.ID); err != nil {
			return err
		}

		if err := s.addGatesStatusTask(ctx, device); err != nil {
			return err
		}
	}

	return s.devicesRepository.UpdateDevice(ctx, device)
}

func (s *Service) addGatesStatusTask(ctx context.Context, device entities.Device) error {
	checkGatesStatusTask := cron.NewCheckGatesStatusTask(s.devicesRepository, s.gatesClient, s.telemetryRepository, device.ID)

	if err := s.cronScheduler.AddJob(ctx, device.Schedule, device.ID, func(ctx context.Context) error {
		return checkGatesStatusTask.Run(ctx)
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteDevice(ctx context.Context, deviceID uuid.UUID) error {
	device, err := s.devicesRepository.GetDevice(ctx, deviceID)
	if err != nil {
		return err
	}

	err = s.gatesClient.DeactivateGate(ctx, device.Host)
	if err != nil {
		return err
	}

	return s.devicesRepository.DeleteDevice(ctx, deviceID)
}

func (s *Service) DeleteHouseDevices(ctx context.Context, houseID uuid.UUID) error {
	devices, err := s.devicesRepository.GetHouseDevices(ctx, houseID)
	if err != nil {
		return err
	}

	for _, device := range devices {
		if err := s.DeleteDevice(ctx, device.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) DeleteUserDevices(ctx context.Context, userID uuid.UUID) error {
	devices, err := s.devicesRepository.GetUserDevices(ctx, userID)
	if err != nil {
		return err
	}

	for _, device := range devices {
		if err := s.DeleteDevice(ctx, device.ID); err != nil {
			return err
		}
	}

	return nil
}
