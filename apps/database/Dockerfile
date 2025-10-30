FROM alpine:latest

# Устанавливаем необходимые пакеты
RUN apk add --no-cache \
    curl \
    postgresql-client

# Копируем скрипт миграций
COPY run-migrations.sh /usr/local/bin/run-migrations.sh
RUN chmod 755 /usr/local/bin/run-migrations.sh

# Копируем миграции
COPY migrations /migrations

# Запускаем миграции
CMD ["/bin/sh", "/usr/local/bin/run-migrations.sh"]
