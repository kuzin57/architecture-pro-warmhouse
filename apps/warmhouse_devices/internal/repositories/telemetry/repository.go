package telemetry

import (
	"context"
	"fmt"

	"github.com/warmhouse/warmhouse_devices/internal/entities"
	"github.com/warmhouse/warmhouse_devices/internal/repositories"
)

type Repository struct {
	driver *repositories.PgDriver
}

func NewRepository(driver *repositories.PgDriver) *Repository {
	return &Repository{driver: driver}
}

func (r *Repository) AddTelemetry(ctx context.Context, telemetry entities.Telemetry) error {
	query := `
		INSERT INTO warmhouse.telemetry (device_id, timestamp, data)
		VALUES (:device_id, :timestamp, :data)
	`

	_, err := r.driver.DB().NamedExecContext(ctx, query, telemetry)
	if err != nil {
		return fmt.Errorf("error adding telemetry: %w", err)
	}

	return nil
}
