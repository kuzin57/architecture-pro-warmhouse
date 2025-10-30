package common

import (
	"context"

	"github.com/google/uuid"
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
	GetAllActiveDevices(ctx context.Context) ([]entities.Device, error)
}
