package router

import (
	"context"

	"github.com/google/uuid"
	"github.com/warmhouse/warmhouse_devices/internal/generated/async"
)

type DevicesService interface {
	CreateDevice(ctx context.Context, payload async.DeviceCreatedMessageFromDeviceCreatedChannelPayload) error
	UpdateDevice(ctx context.Context, payload async.DeviceUpdatedMessageFromDeviceUpdatedChannelPayload) error
	DeleteDevice(ctx context.Context, deviceID uuid.UUID) error
	DeleteHouseDevices(ctx context.Context, houseID uuid.UUID) error
	DeleteUserDevices(ctx context.Context, userID uuid.UUID) error
}

type WarmingService interface {
	DevicesService
}

type GatesService interface {
	DevicesService
}
