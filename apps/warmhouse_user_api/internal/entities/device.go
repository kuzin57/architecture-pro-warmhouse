package entities

import (
	"time"

	"github.com/warmhouse/warmhouse_user_api/internal/consts"

	"github.com/google/uuid"
)

type Device struct {
	ID        uuid.UUID           `db:"id"`
	HouseID   uuid.UUID           `db:"house_id"`
	Name      string              `db:"name"`
	Location  string              `db:"location"`
	Host      string              `db:"host"`
	Unit      string              `db:"unit"`
	Value     string              `db:"value"`
	SensorID  *int64              `db:"sensor_id"`
	Schedule  string              `db:"schedule"`
	Status    consts.DeviceStatus `db:"status"`
	Type      consts.DeviceType   `db:"type"`
	CreatedAt time.Time           `db:"created_at"`
	UpdatedAt time.Time           `db:"updated_at"`
}
