package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/warmhouse/warmhouse_devices/internal/consts"
)

type Device struct {
	ID        uuid.UUID           `db:"id"`
	HouseID   uuid.UUID           `db:"house_id"`
	Name      string              `db:"name"`
	SensorID  *int64              `db:"sensor_id"`
	Location  string              `db:"location"`
	Schedule  string              `db:"schedule"`
	Host      string              `db:"host"`
	Unit      string              `db:"unit"`
	Value     string              `db:"value"`
	Status    consts.DeviceStatus `db:"status"`
	Type      consts.DeviceType   `db:"type"`
	CreatedAt time.Time           `db:"created_at"`
	UpdatedAt time.Time           `db:"updated_at"`
}

func GetDeviceStatuses() []string {
	return []string{
		string(consts.DeviceStatusInactive),
		string(consts.DeviceStatusActive),
		string(consts.DeviceStatusUnknown),
	}
}
