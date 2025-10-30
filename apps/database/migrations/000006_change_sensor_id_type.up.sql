ALTER TABLE warmhouse.devices
DROP COLUMN IF EXISTS sensor_id;

ALTER TABLE warmhouse.devices
ADD COLUMN IF NOT EXISTS sensor_id INTEGER;

COMMENT ON COLUMN warmhouse.devices.sensor_id IS 'ID сенсора устройства (например: 1)';
