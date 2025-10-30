package cron

import (
	"context"

	"github.com/google/uuid"
	gatesapi "github.com/warmhouse/warmhouse_devices/internal/clients/gates_api"
	temperatureapi "github.com/warmhouse/warmhouse_devices/internal/clients/temperature_api"
	"github.com/warmhouse/warmhouse_devices/internal/consts"
	"github.com/warmhouse/warmhouse_devices/internal/entities"
)

type DevicesRepository interface {
	GetDevice(ctx context.Context, deviceID uuid.UUID) (entities.Device, error)
	SetDeviceStatus(ctx context.Context, deviceID uuid.UUID, status consts.DeviceStatus) error
}

type GatesClient interface {
	GetGateStatus(ctx context.Context, host string) (*gatesapi.StatusResponse, error)
}

type TemperatureClient interface {
	GetTemperature(ctx context.Context, host string, location string) (*temperatureapi.TemperatureResponse, error)
	GetTemperatureBySensorID(ctx context.Context, host string, sensorID string) (*temperatureapi.TemperatureResponse, error)
	HealthCheck(ctx context.Context, host string) error
}

type TelemetryRepository interface {
	AddTelemetry(ctx context.Context, telemetry entities.Telemetry) error
}
