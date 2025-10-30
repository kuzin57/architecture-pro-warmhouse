package gates

import (
	"context"

	"github.com/google/uuid"
	gatesapi "github.com/warmhouse/warmhouse_devices/internal/clients/gates_api"
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

type GatesClient interface {
	ActivateGate(ctx context.Context, host string) error
	DeactivateGate(ctx context.Context, host string) error
	GetGateStatus(ctx context.Context, host string) (*gatesapi.StatusResponse, error)
	ChangeGateState(ctx context.Context, host string, state gatesapi.GateState) error
	HealthCheck(ctx context.Context, host string) error
}

type CronScheduler interface {
	AddJob(ctx context.Context, schedule string, job func(ctx context.Context) error) error
}
