package houses

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

func (r *Repository) CreateHouse(ctx context.Context, house entities.House) error {
	query := `
		INSERT INTO warmhouse.houses (id, user_id, name, address)
		VALUES (:id, :user_id, :name, :address)
	`

	_, err := r.driver.DB().NamedExecContext(ctx, query, house)
	if err != nil {
		return fmt.Errorf("error creating house: %w", err)
	}

	return nil
}

func (r *Repository) GetHouse(ctx context.Context, houseID, userID uuid.UUID) (entities.House, error) {
	query := `
		SELECT * FROM warmhouse.houses WHERE id = $1 AND user_id = $2
	`

	var house entities.House
	err := r.driver.DB().GetContext(ctx, &house, query, houseID, userID)
	if err != nil {
		return entities.House{}, fmt.Errorf("error getting house: %w", err)
	}

	return house, nil
}

func (r *Repository) GetUserHouses(ctx context.Context, userID uuid.UUID) ([]entities.House, error) {
	query := `
		SELECT * FROM warmhouse.houses WHERE user_id = $1
	`

	var houses []entities.House
	err := r.driver.DB().SelectContext(ctx, &houses, query, userID)
	if err != nil {
		return []entities.House{}, fmt.Errorf("error getting user houses: %w", err)
	}

	return houses, nil
}

func (r *Repository) UpdateHouse(ctx context.Context, house entities.House) error {
	query := `
		UPDATE warmhouse.houses SET name = :name, address = :address WHERE id = :id
	`

	_, err := r.driver.DB().NamedExecContext(ctx, query, house)
	if err != nil {
		return fmt.Errorf("error updating house: %w", err)
	}

	return nil
}

func (r *Repository) DeleteHouse(ctx context.Context, houseID uuid.UUID) error {
	query := `
		DELETE FROM warmhouse.houses WHERE id = $1
	`

	_, err := r.driver.DB().ExecContext(ctx, query, houseID)
	if err != nil {
		return fmt.Errorf("error deleting house: %w", err)
	}

	return nil
}

func (r *Repository) GetOldestUserHouse(ctx context.Context, userID uuid.UUID) (entities.House, error) {
	query := `
		SELECT * FROM warmhouse.houses 
		WHERE user_id = $1 
		ORDER BY created_at ASC 
		LIMIT 1
	`

	var house entities.House
	err := r.driver.DB().GetContext(ctx, &house, query, userID)
	if err != nil {
		return entities.House{}, fmt.Errorf("error getting oldest user house: %w", err)
	}

	return house, nil
}
