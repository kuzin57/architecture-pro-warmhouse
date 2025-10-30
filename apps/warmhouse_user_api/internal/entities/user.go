package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID      `db:"id"`
	Email          string         `db:"email"`
	Phone          sql.NullString `db:"phone"`
	Name           string         `db:"name"`
	HashedPassword string         `db:"hashed_password"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
}
