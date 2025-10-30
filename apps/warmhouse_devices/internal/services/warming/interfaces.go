package warming

import (
	"context"

	"github.com/google/uuid"
	temperatureapi "github.com/warmhouse/warmhouse_devices/internal/clients/temperature_api"
	"github.com/warmhouse/warmhouse_devices/internal/consts"
	"github.com/warmhouse/warmhouse_devices/internal/entities"
)

type DevicesRepository interface {
	CreateDevice(ctx context.Context, device entities.Device) error
	GetDevice(ctx context.Context, deviceID uuid.UUID) (entities.Device, error)
	UpdateDevice(ctx context.Context, device entities.Device) error
	DeleteDevice(ctx context.Context, deviceID uuid.UUID) error
	GetHouseDevices(ctx context.Context, houseID uuid.UUID) ([]entities.Device, error)
	GetUserDevices(ctx context.Context, userID uuid.UUID) ([]entities.Device, error)
	SetDeviceStatus(ctx context.Context, deviceID uuid.UUID, status consts.DeviceStatus) error
}

type TemperatureClient interface {
	GetTemperature(ctx context.Context, host string, location string) (*temperatureapi.TemperatureResponse, error)
	GetTemperatureBySensorID(ctx context.Context, host string, sensorID string) (*temperatureapi.TemperatureResponse, error)
	HealthCheck(ctx context.Context, host string) error
}

type CronScheduler interface {
	AddJob(ctx context.Context, schedule string, job func(ctx context.Context) error) error
}
