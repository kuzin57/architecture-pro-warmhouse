package warming

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/warmhouse/libraries/convert"
	"github.com/warmhouse/libraries/scheduler"
	"github.com/warmhouse/warmhouse_devices/internal/consts"
	"github.com/warmhouse/warmhouse_devices/internal/cron"
	"github.com/warmhouse/warmhouse_devices/internal/entities"
	"github.com/warmhouse/warmhouse_devices/internal/generated/async"
	"github.com/warmhouse/warmhouse_devices/internal/repositories/telemetry"
)

type Service struct {
	devicesRepository   DevicesRepository
	telemetryRepository *telemetry.Repository
	temperatureClient   TemperatureClient
	cronScheduler       *scheduler.Scheduler
}

func NewService(
	devicesRepository DevicesRepository,
	telemetryRepository *telemetry.Repository,
	temperatureClient TemperatureClient,
	cronScheduler *scheduler.Scheduler,
) *Service {
	return &Service{
		devicesRepository:   devicesRepository,
		telemetryRepository: telemetryRepository,
		temperatureClient:   temperatureClient,
		cronScheduler:       cronScheduler,
	}
}

func (s *Service) CreateDevice(ctx context.Context, payload async.DeviceCreatedMessageFromDeviceCreatedChannelPayload) error {
	var (
		device = entities.Device{
			ID:       uuid.MustParse(payload.DeviceId),
			HouseID:  uuid.MustParse(payload.HouseId),
			Name:     payload.Name,
			Location: convert.UnwrapOr(payload.Location, ""),
			SensorID: payload.SensorId,
			Host:     payload.Host,
			Unit:     consts.DeviceUnitWarming,
			Status:   consts.DeviceStatusUnknown,
			Type:     consts.DeviceTypeTemperature,
			Value:    convert.UnwrapOr(payload.Value, "-"),
			Schedule: convert.UnwrapOr(payload.Schedule, checkTemperatureDefaultSchedule),
		}
		checkTemperatureTask = cron.NewCheckTemperatureTask(s.devicesRepository, s.temperatureClient, s.telemetryRepository, device.ID)
	)

	if err := s.cronScheduler.AddJob(ctx, device.Schedule, device.ID, func(ctx context.Context) error {
		if err := checkTemperatureTask.Run(ctx); err != nil {
			log.Println("error running check temperature task", err)

			return err
		}

		return nil
	}); err != nil {
		return err
	}

	if err := s.devicesRepository.CreateDevice(ctx, device); err != nil {
		log.Println("error creating device", err)

		return err
	}

	return nil
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
	device.Location = convert.UnwrapOr(payload.Location, device.Location)
	device.SensorID = payload.SensorId
	device.Host = convert.UnwrapOr(payload.Host, device.Host)
	device.Unit = convert.UnwrapOr(payload.Unit, device.Unit)
	device.Value = convert.UnwrapOr(payload.Value, device.Value)
	device.Status = convert.UnwrapOr((*consts.DeviceStatus)(payload.Status), device.Status)

	switch {
	case (payload.Status != nil && *payload.Status != string(device.Status) && *payload.Status == string(consts.DeviceStatusActive)) ||
		(payload.Schedule != nil && *payload.Schedule != device.Schedule):
		if err := s.cronScheduler.RemoveJob(ctx, device.ID); err != nil {
			return err
		}

		fallthrough
	case (payload.Status != nil && *payload.Status != string(device.Status) && *payload.Status == string(consts.DeviceStatusInactive)) ||
		(payload.Schedule != nil && *payload.Schedule != device.Schedule):
		checkTemperatureTask := cron.NewCheckTemperatureTask(s.devicesRepository, s.temperatureClient, s.telemetryRepository, device.ID)

		if err := s.cronScheduler.AddJob(ctx, device.Schedule, device.ID, func(ctx context.Context) error {
			return checkTemperatureTask.Run(ctx)
		}); err != nil {
			return err
		}
	}

	return s.devicesRepository.UpdateDevice(ctx, device)
}

func (s *Service) DeleteDevice(ctx context.Context, deviceID uuid.UUID) error {
	if err := s.cronScheduler.RemoveJob(ctx, deviceID); err != nil {
		log.Println("error removing cron job", err)

		return err
	}

	if err := s.devicesRepository.DeleteDevice(ctx, deviceID); err != nil {
		log.Println("error deleting device", err)

		return err
	}

	return nil
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
