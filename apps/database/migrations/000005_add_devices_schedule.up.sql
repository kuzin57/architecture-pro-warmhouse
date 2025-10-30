ALTER TABLE warmhouse.devices
ADD COLUMN IF NOT EXISTS schedule TEXT NOT NULL DEFAULT '* * * * *';

COMMENT ON COLUMN warmhouse.devices.schedule IS 'Расписание устройства (например: "0 0 * * *", "0 0 * * * *")';
