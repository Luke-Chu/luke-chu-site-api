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
- `POST /api/v1/photos/:uuid/view` (implemented with anti-abuse window)
- `POST /api/v1/photos/:uuid/like` (implemented with dedup like)
- `POST /api/v1/photos/:uuid/unlike` (implemented with count rollback)
- `POST /api/v1/photos/:uuid/download` (implemented with anti-abuse window + url)
- `GET /api/v1/tags`
- `GET /api/v1/filters`

## New In This Step

- Added anti-abuse window for `POST /api/v1/photos/:uuid/download`.
- Download anti-abuse behavior:
  - same visitor within 30 minutes is not counted repeatedly
  - records download in `photo_downloads`
  - response includes `counted` to indicate whether this request increased count

## Build

```bash
./scripts/build.sh
```
