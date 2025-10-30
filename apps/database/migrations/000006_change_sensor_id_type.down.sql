ALTER TABLE warmhouse.devices
DROP COLUMN IF EXISTS sensor_id;

ALTER TABLE warmhouse.devices
ADD COLUMN IF NOT EXISTS sensor_id TEXT NOT NULL DEFAULT '';

COMMENT ON COLUMN warmhouse.devices.sensor_id IS 'ID сенсора устройства (например: "1234567890")';
