package users

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

func (r *Repository) CreateUser(ctx context.Context, user entities.User) error {
	query := `
		INSERT INTO warmhouse.users (id, email, phone, name, hashed_password)
		VALUES (:id, :email, :phone, :name, :hashed_password)
	`

	_, err := r.driver.DB().NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (r *Repository) GetUser(ctx context.Context, userID uuid.UUID) (entities.User, error) {
	query := `
		SELECT * FROM warmhouse.users WHERE id = $1
	`

	var user entities.User
	err := r.driver.DB().GetContext(ctx, &user, query, userID)
	if err != nil {
		return entities.User{}, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user entities.User) error {
	query := `
		UPDATE warmhouse.users
		SET email = :email,
			phone = :phone,
			name = :name,
			hashed_password = :hashed_password
		WHERE id = :id
	`

	_, err := r.driver.DB().NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (entities.User, error) {
	query := `
		SELECT * FROM warmhouse.users WHERE email = $1
	`

	var user entities.User
	err := r.driver.DB().GetContext(ctx, &user, query, email)
	if err != nil {
		return entities.User{}, fmt.Errorf("error getting user by email: %w", err)
	}

	return user, nil
}

func (r *Repository) GetDefaultUser(ctx context.Context) (entities.User, error) {
	query := `
		SELECT *
		FROM warmhouse.users
		ORDER BY created_at ASC
		LIMIT 1
	`

	var user entities.User
	err := r.driver.DB().GetContext(ctx, &user, query)
	if err != nil {
		return entities.User{}, fmt.Errorf("error getting default user: %w", err)
	}

	return user, nil
}
