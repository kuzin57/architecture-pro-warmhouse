package telemetry

import (
	"context"
	"fmt"
	"time"

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

func (r *Repository) GetDeviceTelemetry(ctx context.Context, deviceID uuid.UUID, limit, offset int, fromDate, toDate *time.Time) ([]entities.Telemetry, error) {
	query := `
		SELECT device_id, timestamp, data
		FROM warmhouse.telemetry
		WHERE device_id = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3
	`

	var telemetry []entities.Telemetry
	err := r.driver.DB().SelectContext(ctx, &telemetry, query, deviceID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting device telemetry: %w", err)
	}

	return telemetry, nil
}
