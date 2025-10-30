-- Откат миграции добавления location и host в таблицу devices

-- Переключение на базу данных warmhouse
\c warmhouse

-- Удаление индексов
DROP INDEX IF EXISTS idx_devices_host;
DROP INDEX IF EXISTS idx_devices_location;

-- Удаление колонок
ALTER TABLE warmhouse.devices DROP COLUMN IF EXISTS host;
ALTER TABLE warmhouse.devices DROP COLUMN IF EXISTS location;
