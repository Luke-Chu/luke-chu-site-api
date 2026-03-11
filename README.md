# luke-chu-site-api

Backend service for `luke-chu.com`. Current focus is the photography module.

## Tech Stack

- Go
- Gin
- sqlx
- PostgreSQL
- Viper (YAML config)
- Zap
- go-playground/validator
- google/uuid

## Run Local

```bash
APP_ENV=local go run ./cmd/server
```

Default listen address: `:8080`.

## Database (local default)

- host: `127.0.0.1`
- port: `5432`
- user: `admin`
- password: `1234`
- dbname: `luke-chu-site`

## Current APIs

- `GET /api/v1/health`
- `GET /api/v1/photos`
- `GET /api/v1/photos/:uuid` (implemented with full detail + tags)
- `POST /api/v1/photos/:uuid/view` (implemented with count increment)
- `POST /api/v1/photos/:uuid/like` (skeleton, to be finalized)
- `POST /api/v1/photos/:uuid/download` (skeleton, to be finalized)
- `GET /api/v1/tags`
- `GET /api/v1/filters`

## New In This Step

- Completed `GET /api/v1/photos/:uuid` and `POST /api/v1/photos/:uuid/view`.
- View API behavior:
  - UUID validation
  - `is_published = true` update with `view_count = view_count + 1`
  - returns latest `viewCount`
  - `400 invalid uuid`, `404 photo not found`, `500 internal server error`

## Build

```bash
./scripts/build.sh
```
