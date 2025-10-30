package devices

import (
	"context"
	"fmt"

	"github.com/warmhouse/warmhouse_user_api/internal/entities"
	"github.com/warmhouse/warmhouse_user_api/internal/repositories"

	"github.com/google/uuid"
)

type Repository struct {
	driver *repositories.PgDriver
}

func NewRepository(driver *repositories.PgDriver) *Repository {
	return &Repository{driver: driver}
}

func (r *Repository) CreateDevice(ctx context.Context, device entities.Device) error {
	query := `
		INSERT INTO warmhouse.devices (id, house_id, name, location, host, unit, value, status)
		VALUES (:id, :house_id, :name, :location, :host, :unit, :value, :status)
	`

	_, err := r.driver.DB().NamedExecContext(ctx, query, device)
	if err != nil {
		return fmt.Errorf("error creating device: %w", err)
	}

	return nil
}

func (r *Repository) GetDevice(ctx context.Context, deviceID, userID uuid.UUID) (entities.Device, error) {
	query := `
		SELECT 
			warmhouse.devices.*
		FROM warmhouse.devices
		JOIN warmhouse.houses ON warmhouse.devices.house_id = warmhouse.houses.id
		WHERE warmhouse.devices.id = $1 AND warmhouse.houses.user_id = $2
	`

	var device entities.Device
	err := r.driver.DB().GetContext(ctx, &device, query, deviceID, userID)
	if err != nil {
		return entities.Device{}, fmt.Errorf("error getting device info: %w", err)
	}

	return device, nil
}

func (r *Repository) GetUserDevices(ctx context.Context, userID uuid.UUID) ([]entities.Device, error) {
	query := `
		SELECT warmhouse.devices.*
		FROM warmhouse.devices
		JOIN warmhouse.houses ON warmhouse.devices.house_id = warmhouse.houses.id
		WHERE warmhouse.houses.user_id = $1
	`

	var devices []entities.Device
	err := r.driver.DB().SelectContext(ctx, &devices, query, userID)
	if err != nil {
		return []entities.Device{}, fmt.Errorf("error getting user devices: %w", err)
	}

	return devices, nil
}

func (r *Repository) GetHouseDevices(ctx context.Context, houseID, userID uuid.UUID) ([]entities.Device, error) {
	query := `
		SELECT warmhouse.devices.*
		FROM warmhouse.devices
		JOIN warmhouse.houses ON warmhouse.devices.house_id = warmhouse.houses.id
		WHERE warmhouse.houses.user_id = $1 AND warmhouse.devices.house_id = $2
	`

	var devices []entities.Device
	err := r.driver.DB().SelectContext(ctx, &devices, query, userID, houseID)
	if err != nil {
		return []entities.Device{}, fmt.Errorf("error getting house devices: %w", err)
	}

	return devices, nil
}
