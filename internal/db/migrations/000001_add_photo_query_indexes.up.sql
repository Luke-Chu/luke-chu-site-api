-- Photo list / filters query performance baseline
-- Safe to run multiple times.

CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Sort-friendly partial indexes for published photos
CREATE INDEX IF NOT EXISTS idx_photos_pub_shot_time
ON photos (shot_time DESC, id DESC)
WHERE is_published = TRUE;

CREATE INDEX IF NOT EXISTS idx_photos_pub_like_count
ON photos (like_count DESC, id DESC)
WHERE is_published = TRUE;

CREATE INDEX IF NOT EXISTS idx_photos_pub_view_count
ON photos (view_count DESC, id DESC)
WHERE is_published = TRUE;

CREATE INDEX IF NOT EXISTS idx_photos_pub_download_count
ON photos (download_count DESC, id DESC)
WHERE is_published = TRUE;

CREATE INDEX IF NOT EXISTS idx_photos_pub_created_at
ON photos (created_at DESC, id DESC)
WHERE is_published = TRUE;

-- Common filter indexes
CREATE INDEX IF NOT EXISTS idx_photos_pub_year_month
ON photos (year DESC, month DESC, id DESC)
WHERE is_published = TRUE;

CREATE INDEX IF NOT EXISTS idx_photos_pub_orientation
ON photos (orientation, id DESC)
WHERE is_published = TRUE;

CREATE INDEX IF NOT EXISTS idx_photos_pub_category
ON photos (category, id DESC)
WHERE is_published = TRUE;

-- Keyword search indexes (ILIKE on lower(...))
CREATE INDEX IF NOT EXISTS idx_photos_title_cn_trgm
ON photos USING gin (lower(COALESCE(title_cn, '')) gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_photos_title_en_trgm
ON photos USING gin (lower(COALESCE(title_en, '')) gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_photos_filename_trgm
ON photos USING gin (lower(filename) gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_tags_name_lower
ON tags (lower(name));

CREATE INDEX IF NOT EXISTS idx_tags_name_trgm
ON tags USING gin (lower(name) gin_trgm_ops);

-- Join acceleration for tag filters
CREATE INDEX IF NOT EXISTS idx_photo_tags_tag_id_photo_id
ON photo_tags (tag_id, photo_id);

