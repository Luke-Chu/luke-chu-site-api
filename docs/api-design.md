# API Design

## Base Path

- `/api/v1`

## Response Envelope

Success:

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

Error:

```json
{
  "code": 40002,
  "message": "invalid uuid",
  "data": null
}
```

## Implemented Routes

- `GET /api/v1/health`
- `GET /api/v1/photos`
- `GET /api/v1/photos/:uuid`
- `POST /api/v1/photos/:uuid/view`
- `POST /api/v1/photos/:uuid/like` (to be finalized)
- `POST /api/v1/photos/:uuid/download` (to be finalized)
- `GET /api/v1/tags`
- `GET /api/v1/filters`

## GET /api/v1/photos/:uuid

### Purpose

Return published photo detail for the detail page.

### Path Param

- `uuid`: photo UUID

### Behavior

- Validate UUID format.
- Query photo with `is_published = true`.
- Query all tags for this photo.
- Return detail fields + tags.

### Error Semantics

- `400`: invalid UUID
- `404`: photo not found
- `500`: database/internal error

### Response Fields (core)

- `id`, `uuid`, `filename`
- `titleCn`, `titleEn`, `description`, `category`
- `shotTime`, `width`, `height`, `resolution`, `orientation`
- `cameraModel`, `lensModel`, `focalLength`, `focalLength35mm`
- `aperture`, `shutterSpeed`, `iso`
- `meteringMode`, `exposureCompensation`, `exposureProgram`
- `whiteBalance`, `flash`
- `thumbUrl`, `displayUrl`, `originalUrl`
- `likeCount`, `viewCount`, `downloadCount`
- `createdAt`, `updatedAt`
- `tags` (`id`, `name`, `tagType`)

## POST /api/v1/photos/:uuid/view

### Purpose

Increase view count and return latest `viewCount`.

### Path Param

- `uuid`: photo UUID

### Behavior

- Validate UUID format.
- Update published photo only:
  - `view_count = view_count + 1`
  - `updated_at = NOW()`
- Return `uuid` and latest `viewCount`.

### Error Semantics

- `400`: invalid UUID
- `404`: photo not found
- `500`: database/internal error
