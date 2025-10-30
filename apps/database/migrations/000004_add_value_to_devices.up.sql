ALTER TABLE warmhouse.devices
ADD COLUMN IF NOT EXISTS value TEXT NOT NULL DEFAULT '';

COMMENT ON COLUMN warmhouse.devices.value IS 'Значение устройства (например: "1234567890")';
