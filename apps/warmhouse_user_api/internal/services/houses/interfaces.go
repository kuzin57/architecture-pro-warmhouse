package houses

import (
	"context"

	"github.com/warmhouse/warmhouse_user_api/internal/entities"

	"github.com/google/uuid"
)

type HousesRepository interface {
	GetUserHouses(ctx context.Context, userID uuid.UUID) ([]entities.House, error)
	CreateHouse(ctx context.Context, house entities.House) error
	GetHouse(ctx context.Context, houseID, userID uuid.UUID) (entities.House, error)
	UpdateHouse(ctx context.Context, house entities.House) error
	DeleteHouse(ctx context.Context, houseID uuid.UUID) error
	GetOldestUserHouse(ctx context.Context, userID uuid.UUID) (entities.House, error)
}

type DevicesRepository interface {
	GetHouseDevices(ctx context.Context, houseID, userID uuid.UUID) ([]entities.Device, error)
}
