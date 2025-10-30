package common

import (
	"context"
	"log"

	"github.com/warmhouse/libraries/scheduler"
	gatesapi "github.com/warmhouse/warmhouse_devices/internal/clients/gates_api"
	temperatureapi "github.com/warmhouse/warmhouse_devices/internal/clients/temperature_api"
	"github.com/warmhouse/warmhouse_devices/internal/consts"
	"github.com/warmhouse/warmhouse_devices/internal/cron"
	"github.com/warmhouse/warmhouse_devices/internal/repositories/devices"
	"github.com/warmhouse/warmhouse_devices/internal/repositories/telemetry"
	"go.uber.org/fx"
)

type Starter struct {
	devicesRepository   *devices.Repository
	gatesClient         *gatesapi.Client
	temperatureClient   *temperatureapi.Client
	telemetryRepository *telemetry.Repository
	cronScheduler       *scheduler.Scheduler
}

func NewStarter(
	lc fx.Lifecycle,
	devicesRepository *devices.Repository,
	gatesClient *gatesapi.Client,
	telemetryRepository *telemetry.Repository,
	temperatureClient *temperatureapi.Client,
	cronScheduler *scheduler.Scheduler,
) *Starter {
	s := &Starter{
		devicesRepository:   devicesRepository,
		gatesClient:         gatesClient,
		telemetryRepository: telemetryRepository,
		temperatureClient:   temperatureClient,
		cronScheduler:       cronScheduler,
	}

	lc.Append(fx.Hook{
		OnStart: s.Start,
		OnStop:  s.Stop,
	})

	return s
}

func (s *Starter) Start(ctx context.Context) error {
	devices, err := s.devicesRepository.GetAllActiveDevices(ctx)
	if err != nil {
		return err
	}

	for _, device := range devices {
		log.Println("starting device healthcheck", device.ID)

		switch device.Type {
		case consts.DeviceTypeTemperature:
			checkTemperatureTask := cron.NewCheckTemperatureTask(s.devicesRepository, s.temperatureClient, s.telemetryRepository, device.ID)
			if err := s.cronScheduler.AddJob(ctx, device.Schedule, device.ID, func(ctx context.Context) error {
				return checkTemperatureTask.Run(ctx)
			}); err != nil {
				return err
			}
		case consts.DeviceTypeGates:
			checkGatesStatusTask := cron.NewCheckGatesStatusTask(s.devicesRepository, s.gatesClient, s.telemetryRepository, device.ID)
			if err := s.cronScheduler.AddJob(ctx, device.Schedule, device.ID, func(ctx context.Context) error {
				return checkGatesStatusTask.Run(ctx)
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Starter) Stop(ctx context.Context) error {
	return nil
}
