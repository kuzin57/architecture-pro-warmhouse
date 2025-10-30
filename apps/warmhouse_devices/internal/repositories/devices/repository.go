package devices

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/warmhouse/warmhouse_devices/internal/consts"
	"github.com/warmhouse/warmhouse_devices/internal/entities"
	"github.com/warmhouse/warmhouse_devices/internal/repositories"
)

type Repository struct {
	driver *repositories.PgDriver
}

func NewRepository(driver *repositories.PgDriver) *Repository {
	return &Repository{driver: driver}
}

func (r *Repository) CreateDevice(ctx context.Context, device entities.Device) error {
	query := `
		INSERT INTO warmhouse.devices
			(id, house_id, name, sensor_id, unit, value, status, type, host, location, schedule)
		VALUES
			(:id, :house_id, :name, :sensor_id, :unit, :value, :status, :type, :host, :location, :schedule)
	`

	_, err := r.driver.DB().NamedExecContext(ctx, query, device)
	if err != nil {
		return fmt.Errorf("error creating device: %w", err)
	}

	return nil
}

func (r *Repository) GetDevice(ctx context.Context, deviceID uuid.UUID) (entities.Device, error) {
	query := `
		SELECT * FROM warmhouse.devices WHERE id = $1
	`

	var device entities.Device
	err := r.driver.DB().GetContext(ctx, &device, query, deviceID)
	if err != nil {
		return entities.Device{}, fmt.Errorf("error getting device: %w", err)
	}

	return device, nil
}

func (r *Repository) UpdateDevice(ctx context.Context, device entities.Device) error {
	query := `
		UPDATE warmhouse.devices
		SET name = :name,
			sensor_id = COALESCE(:sensor_id, sensor_id),
			unit = :unit,
			value = :value,
			status = :status,
			location = :location,
			host = :host,
			schedule = :schedule
		WHERE id = :id
	`

	_, err := r.driver.DB().NamedExecContext(ctx, query, device)
	if err != nil {
		return fmt.Errorf("error updating device: %w", err)
	}

	return nil
}

func (r *Repository) DeleteDevice(ctx context.Context, deviceID uuid.UUID) error {
	query := `
		DELETE FROM warmhouse.devices
		WHERE id = $1
	`

	_, err := r.driver.DB().ExecContext(ctx, query, deviceID)
	if err != nil {
		return fmt.Errorf("error deleting device: %w", err)
	}

	return nil
}

func (r *Repository) GetHouseDevices(ctx context.Context, houseID uuid.UUID) ([]entities.Device, error) {
	query := `
		SELECT * FROM warmhouse.devices WHERE house_id = $1
	`

	var devices []entities.Device
	err := r.driver.DB().SelectContext(ctx, &devices, query, houseID.String())
	if err != nil {
		return nil, fmt.Errorf("error getting house devices: %w", err)
	}

	return devices, nil
}

func (r *Repository) GetUserDevices(ctx context.Context, userID uuid.UUID) ([]entities.Device, error) {
	query := `
		SELECT *
		FROM warmhouse.devices
		JOIN warmhouse.houses ON warmhouse.devices.house_id = warmhouse.houses.id
		WHERE warmhouse.houses.user_id = $1
	`

	var devices []entities.Device
	err := r.driver.DB().SelectContext(ctx, &devices, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user devices: %w", err)
	}

	return devices, nil
}

func (r *Repository) SetDeviceStatus(ctx context.Context, deviceID uuid.UUID, status consts.DeviceStatus) error {
	query := `
		UPDATE warmhouse.devices
		SET status = $1
		WHERE id = $2
	`

	_, err := r.driver.DB().ExecContext(ctx, query, status, deviceID)
	if err != nil {
		return fmt.Errorf("error setting device status: %w", err)
	}

	return nil
}

func (r *Repository) GetAllActiveDevices(ctx context.Context) ([]entities.Device, error) {
	query := `
		SELECT *
		FROM warmhouse.devices
		WHERE status = 'active' or status = 'unknown'
	`

	var devices []entities.Device
	err := r.driver.DB().SelectContext(ctx, &devices, query)
	if err != nil {
		return []entities.Device{}, fmt.Errorf("error getting all active devices: %w", err)
	}

	return devices, nil
}
