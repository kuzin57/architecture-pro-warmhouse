package devices

import (
	"context"
	"time"

	"github.com/warmhouse/warmhouse_user_api/internal/entities"

	"github.com/google/uuid"
)

type DevicesRepository interface {
	GetDevice(ctx context.Context, deviceID, userID uuid.UUID) (entities.Device, error)
	GetUserDevices(ctx context.Context, userID uuid.UUID) ([]entities.Device, error)
	GetHouseDevices(ctx context.Context, houseID, userID uuid.UUID) ([]entities.Device, error)
	CreateDevice(ctx context.Context, device entities.Device) error
}

type HouseRepository interface {
	GetHouse(ctx context.Context, houseID, userID uuid.UUID) (entities.House, error)
}

type TelemetryRepository interface {
	GetDeviceTelemetry(ctx context.Context, deviceID uuid.UUID, limit, offset int, fromDate, toDate *time.Time) ([]entities.Telemetry, error)
}
