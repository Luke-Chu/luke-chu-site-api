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
- `POST /api/v1/photos/:uuid/like`
- `POST /api/v1/photos/:uuid/unlike`
- `POST /api/v1/photos/:uuid/download`
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

## POST /api/v1/photos/:uuid/like

### Purpose

Like a published photo with visitor-based deduplication.

### Path Param

- `uuid`: photo UUID

### Behavior

- Validate UUID format.
- Read `visitor_hash` from middleware context.
- Transaction flow:
  - locate published photo by UUID
  - `INSERT INTO photo_likes ... ON CONFLICT DO NOTHING`
  - if inserted: `like_count + 1`
  - if conflict: keep count unchanged
- Return:
  - `liked = true` for first like
  - `liked = false` for duplicate like
  - latest `likeCount`

### Error Semantics

- `400`: invalid UUID / missing visitor hash
- `404`: photo not found
- `500`: database/internal error

## POST /api/v1/photos/:uuid/download

### Purpose

Increase download count and return original download URL.

### Path Param

- `uuid`: photo UUID

### Behavior

- Validate UUID format.
- Update published photo only:
  - `download_count = download_count + 1`
  - `updated_at = NOW()`
- Return:
  - `uuid`
  - latest `downloadCount`
  - `downloadUrl` (`original_url`)

### Error Semantics

- `400`: invalid UUID
- `404`: photo not found
- `500`: database/internal error

## POST /api/v1/photos/:uuid/unlike

### Purpose

Cancel a previous like by the same visitor and rollback `like_count`.

### Path Param

- `uuid`: photo UUID

### Behavior

- Validate UUID format.
- Read `visitor_hash` from middleware context.
- Transaction flow:
  - locate published photo by UUID
  - `DELETE FROM photo_likes WHERE photo_id = ? AND visitor_hash = ?`
  - if deleted: decrement `like_count` (floor at 0)
  - if no record deleted: keep count unchanged
- Return:
  - `unliked = true` if rollback happened
  - latest `likeCount`

### Error Semantics

- `400`: invalid UUID / missing visitor hash
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
