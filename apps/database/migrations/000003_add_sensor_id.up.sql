-- Добавление колонки sensor_id в таблицу devices

ALTER TABLE warmhouse.devices
ADD COLUMN IF NOT EXISTS sensor_id TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_devices_sensor_id ON warmhouse.devices(sensor_id);

COMMENT ON COLUMN warmhouse.devices.sensor_id IS 'ID сенсора устройства (например: "1234567890")';
