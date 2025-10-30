-- Откат миграции добавления sensor_id в таблицу devices

-- Переключение на базу данных warmhouse
\c warmhouse

-- Удаление индекса
DROP INDEX IF EXISTS idx_devices_sensor_id;

-- Удаление колонки
ALTER TABLE warmhouse.devices DROP COLUMN IF EXISTS sensor_id;
