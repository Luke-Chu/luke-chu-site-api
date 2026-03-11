DROP INDEX IF EXISTS idx_photo_tags_tag_id_photo_id;

DROP INDEX IF EXISTS idx_tags_name_trgm;
DROP INDEX IF EXISTS idx_tags_name_lower;

DROP INDEX IF EXISTS idx_photos_filename_trgm;
DROP INDEX IF EXISTS idx_photos_title_en_trgm;
DROP INDEX IF EXISTS idx_photos_title_cn_trgm;

DROP INDEX IF EXISTS idx_photos_pub_category;
DROP INDEX IF EXISTS idx_photos_pub_orientation;
DROP INDEX IF EXISTS idx_photos_pub_year_month;

DROP INDEX IF EXISTS idx_photos_pub_created_at;
DROP INDEX IF EXISTS idx_photos_pub_download_count;
DROP INDEX IF EXISTS idx_photos_pub_view_count;
DROP INDEX IF EXISTS idx_photos_pub_like_count;
DROP INDEX IF EXISTS idx_photos_pub_shot_time;
