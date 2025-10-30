package entities

import (
	"time"

	"github.com/google/uuid"
)

type Telemetry struct {
	DeviceID  uuid.UUID `db:"device_id"`
	Timestamp time.Time `db:"timestamp"`
	Data      string    `db:"data"`
}
