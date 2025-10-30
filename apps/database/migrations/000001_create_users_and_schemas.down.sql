-- Откат миграции создания пользователей и схем

-- Удаление прав доступа
REVOKE ALL ON ALL TABLES IN SCHEMA warmhouse FROM warmhouse_user_api;
REVOKE ALL ON ALL TABLES IN SCHEMA warmhouse FROM warmhouse_devices;
REVOKE USAGE, SELECT ON ALL SEQUENCES IN SCHEMA warmhouse FROM warmhouse_user_api;
REVOKE USAGE, SELECT ON ALL SEQUENCES IN SCHEMA warmhouse FROM warmhouse_devices;

-- Удаление триггеров
DROP TRIGGER IF EXISTS update_devices_updated_at ON warmhouse.devices;
DROP TRIGGER IF EXISTS update_houses_updated_at ON warmhouse.houses;
DROP TRIGGER IF EXISTS update_users_updated_at ON warmhouse.users;

-- Удаление функции
DROP FUNCTION IF EXISTS warmhouse.update_updated_at_column();

-- Удаление индексов
DROP INDEX IF EXISTS idx_devices_house_id;
DROP INDEX IF EXISTS idx_houses_user_id;

-- Удаление таблиц
DROP TABLE IF EXISTS warmhouse.telemetry;
DROP TABLE IF EXISTS warmhouse.devices;
DROP TABLE IF EXISTS warmhouse.houses;
DROP TABLE IF EXISTS warmhouse.users;

-- Удаление типов
DROP TYPE IF EXISTS warmhouse.device_type;
DROP TYPE IF EXISTS warmhouse.device_status;

-- Удаление схемы
DROP SCHEMA IF EXISTS warmhouse;

-- Удаление пользователей
DROP USER IF EXISTS warmhouse_devices;
DROP USER IF EXISTS warmhouse_user_api;
