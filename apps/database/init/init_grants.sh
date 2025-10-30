set -e

# Создание пользователей для микросервисов
# Пароли берутся из переменных окружения
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Создание пользователей для микросервисов
    CREATE USER warmhouse_user_api WITH PASSWORD '$USER_API_PASSWORD';
    CREATE USER warmhouse_devices WITH PASSWORD '$DEVICES_PASSWORD';

    -- Настройка прав доступа

    -- Предоставление доступа к схеме warmhouse
    GRANT USAGE ON SCHEMA warmhouse TO warmhouse_user_api;
    GRANT USAGE ON SCHEMA warmhouse TO warmhouse_devices;

    -- warmhouse_user_api: RW доступ к users, houses, devices + RO доступ к telemetry
    GRANT ALL PRIVILEGES ON TABLE warmhouse.users TO warmhouse_user_api;
    GRANT ALL PRIVILEGES ON TABLE warmhouse.houses TO warmhouse_user_api;
    GRANT ALL PRIVILEGES ON TABLE warmhouse.devices TO warmhouse_user_api;
    GRANT SELECT ON TABLE warmhouse.telemetry TO warmhouse_user_api;

    -- warmhouse_devices: RW доступ к devices и telemetry
    GRANT ALL PRIVILEGES ON TABLE warmhouse.devices TO warmhouse_devices;
    GRANT ALL PRIVILEGES ON TABLE warmhouse.telemetry TO warmhouse_devices;
EOSQL
