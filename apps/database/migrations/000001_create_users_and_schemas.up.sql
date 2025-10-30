-- Включение расширения для генерации UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создание пользователей для микросервисов
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'warmhouse_user_api') THEN
        CREATE USER warmhouse_user_api WITH PASSWORD 'user_api_dev_pass_123';
    END IF;
    
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'warmhouse_devices') THEN
        CREATE USER warmhouse_devices WITH PASSWORD 'devices_dev_pass_123';
    END IF;
END
$$;

-- Создание схемы warmhouse
CREATE SCHEMA IF NOT EXISTS warmhouse;

-- Предоставление прав доступа пользователям
GRANT CONNECT ON DATABASE warmhouse TO warmhouse_user_api;
GRANT CONNECT ON DATABASE warmhouse TO warmhouse_devices;
GRANT USAGE ON SCHEMA warmhouse TO warmhouse_user_api;
GRANT USAGE ON SCHEMA warmhouse TO warmhouse_devices;

-- Создание таблиц в схеме warmhouse

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS warmhouse.users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    phone TEXT,
    name TEXT NOT NULL,
    hashed_password TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Таблица домов
CREATE TABLE IF NOT EXISTS warmhouse.houses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES warmhouse.users(id) ON DELETE CASCADE,
    address TEXT NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TYPE warmhouse.device_status AS ENUM ('active', 'inactive', 'unknown');
CREATE TYPE warmhouse.device_type AS ENUM ('temperature', 'gates', 'video');

-- Таблица устройств
CREATE TABLE IF NOT EXISTS warmhouse.devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    house_id UUID NOT NULL REFERENCES warmhouse.houses(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    unit TEXT,
    status warmhouse.device_status NOT NULL DEFAULT 'unknown',
    type warmhouse.device_type NOT NULL DEFAULT 'temperature',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Таблица телеметрии
CREATE TABLE IF NOT EXISTS warmhouse.telemetry (
    device_id UUID NOT NULL REFERENCES warmhouse.devices(id) ON DELETE CASCADE,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    data TEXT NOT NULL,
    PRIMARY KEY (device_id, timestamp)
);

-- Создание индексов для оптимизации
CREATE INDEX IF NOT EXISTS idx_houses_user_id ON warmhouse.houses(user_id);
CREATE INDEX IF NOT EXISTS idx_devices_house_id ON warmhouse.devices(house_id);

-- Создание функции для обновления updated_at
CREATE OR REPLACE FUNCTION warmhouse.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Создание триггеров для автоматического обновления updated_at

-- Триггер для таблицы users
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON warmhouse.users 
    FOR EACH ROW 
    EXECUTE FUNCTION warmhouse.update_updated_at_column();

-- Триггер для таблицы houses
CREATE TRIGGER update_houses_updated_at 
    BEFORE UPDATE ON warmhouse.houses 
    FOR EACH ROW 
    EXECUTE FUNCTION warmhouse.update_updated_at_column();

-- Триггер для таблицы devices
CREATE TRIGGER update_devices_updated_at 
    BEFORE UPDATE ON warmhouse.devices 
    FOR EACH ROW 
    EXECUTE FUNCTION warmhouse.update_updated_at_column();

-- Предоставление прав на таблицы
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA warmhouse TO warmhouse_user_api;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA warmhouse TO warmhouse_devices;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA warmhouse TO warmhouse_user_api;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA warmhouse TO warmhouse_devices;
