-- Добавление колонок location и host в таблицу devices

-- Добавление колонки location в таблицу devices
-- Колонка обязательна (NOT NULL) с пустой строкой по умолчанию для существующих записей
ALTER TABLE warmhouse.devices 
ADD COLUMN IF NOT EXISTS location TEXT NOT NULL DEFAULT '';

-- Добавление колонки host в таблицу devices
-- Колонка обязательна (NOT NULL) с пустой строкой по умолчанию для существующих записей
ALTER TABLE warmhouse.devices 
ADD COLUMN IF NOT EXISTS host TEXT NOT NULL DEFAULT '';

-- Обновление существующих записей, если они есть
-- Устанавливаем значения по умолчанию для существующих записей
UPDATE warmhouse.devices 
SET location = 'Unknown', host = 'localhost' 
WHERE location = '' OR host = '';

-- Создание индексов для оптимизации поиска
CREATE INDEX IF NOT EXISTS idx_devices_location ON warmhouse.devices(location);
CREATE INDEX IF NOT EXISTS idx_devices_host ON warmhouse.devices(host);

-- Комментарии к колонкам
COMMENT ON COLUMN warmhouse.devices.location IS 'Местоположение устройства в доме (например: "kitchen", "living_room", "bedroom")';
COMMENT ON COLUMN warmhouse.devices.host IS 'Хост устройства для подключения (например: "192.168.1.100", "device.local")';
