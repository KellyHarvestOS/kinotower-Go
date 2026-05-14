# Kinotower REST API

Production-like Go REST API для Kinotower на `net/http`, PostgreSQL и Vue frontend.

## Структура

```text
/cmd/api/main.go              # production entrypoint
/internal/config              # .env config loader
/internal/database            # PostgreSQL, migrations, seeds
/internal/middleware          # logger, recover, CORS, auth, timeout, rate limit
/internal/handlers            # HTTP handlers
/internal/services            # бизнес-логика
/internal/repositories        # SQL repository layer
/internal/models              # DTO и response/request models
/internal/routes              # route registration
/internal/utils               # pagination/format helpers
/internal/auth                # JWT и bcrypt
/internal/response            # reusable JSON responses
/internal/validator           # validation layer
/internal/server              # HTTP server + graceful shutdown
/migrations                   # SQL migrations
/seeds                        # demo seed SQL
/web/templates                # Vue frontend
/web/static                   # статические файлы
```

## .env

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=123456
POSTGRES_DB=kinotower
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
DB_SSLMODE=disable
APP_ADDR=:8080
JWT_SECRET=change-this-secret
AUTO_MIGRATE=true
AUTO_SEED=false
```

## Запуск

```bash
go run ./cmd/api
```

Или:

```bash
make run
```

Frontend:

```text
http://localhost:8080
```

API:

```text
http://localhost:8080/api/v1/films
```

## Миграции

Миграции запускаются автоматически при `AUTO_MIGRATE=true`.

Файлы лежат в `/migrations`:

```text
001_create_dictionaries.up.sql
002_create_users.up.sql
003_create_films.up.sql
004_create_reviews.up.sql
005_create_ratings.up.sql
```

## Seeds

Для загрузки demo данных:

```bash
AUTO_SEED=true go run ./cmd/api
```

Или:

```bash
make seed
```

Demo login:

```text
email: ivanov@ivan.kz
password: asdf1234
```

## Middleware

Middleware подключаются в `cmd/api/main.go` через `middleware.Chain`:

```go
appHandler := middleware.Chain(
    router,
    middleware.Recover(log),
    middleware.RequestID,
    middleware.Logger(log),
    middleware.CORS,
    middleware.JSONContent,
    middleware.RateLimit(cfg.RateLimitRPM),
    middleware.Timeout(cfg.RequestTimeout),
)
```

Закрытые routes защищены JWT middleware:

```go
authMiddleware := middleware.Auth(cfg.JWTSecret)
```

Токен передаётся так:

```text
Authorization: Bearer <token>
```
