# WEB-SRV

## Описание

WEB-SRV — это сервер на Go для хранения и управления документами с поддержкой аутентификации, авторизации, загрузки и скачивания файлов. Использует PostgreSQL и Redis.

## Быстрый старт через Docker Compose

1. **Клонируйте репозиторий:**
   ```sh
   git clone https://github.com/DENFNC/web-test.git
   cd WEB-SRV
   ```
2. **Запустите сервисы:**
   ```sh
   docker-compose up --build
   ```
   Это поднимет три контейнера:
   - `websrv-app` — основной Go-сервис (порт 8080)
   - `postgres` — база данных PostgreSQL (порт 5432)
   - `redis` — кэш Redis (порт 6379)

3. **Остановить сервисы:**
   ```sh
   docker-compose down
   ```

## Переменные окружения

Для работы необходимы переменные окружения (можно задать через .env):

```
DATABASE_URL=postgres://admin:admin@postgres:5432/service?sslmode=disable
DATABASE_MAX_CONNS=25
DATABASE_MIN_CONNS=5
DATABASE_MAX_CONN_LIFE_TIME=30m
DATABASE_MAX_CONN_IDLE_TIME=5m
DATABASE_HEALTH_CHECK_PERIOD=1m
APP_URL=http://localhost:8080
```

## Описание Dockerfile

- Сборка бинарника Go в контейнере `golang:1.21-alpine`.
- Копирование бинарника в минимальный образ `alpine`.
- Открывается порт 8080.
- Стартует приложение командой `./websrv`.

## Описание docker-compose.yaml

- **app**: билдит и запускает основной сервис, пробрасывает порт 8080.
- **postgres**: поднимает PostgreSQL с паролем/логином admin, база service, хранит данные в volume.
- **redis**: поднимает Redis для кэширования.
- Все сервисы объединены в сеть `service`.

## Миграции

Перед первым запуском убедитесь, что применены миграции из папки `migrations/` (или используйте автоматическую миграцию, если реализовано).

## Ручной запуск (без Docker)

1. Установите Go 1.21+ и PostgreSQL, Redis локально.
2. Скопируйте `.env.example` в `.env` и настройте переменные.
3. Установите зависимости:
   ```sh
   go mod download
   ```
4. Соберите и запустите:
   ```sh
   go build -o websrv ./cmd/server
   ./websrv
   ```

---