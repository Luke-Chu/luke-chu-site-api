-- 建议在有代表性的测试数据量下执行
-- 目的：验证 /api/v1/photos 对应 SQL 是否命中索引、排序是否走索引、是否出现大范围回表。

EXPLAIN (ANALYZE, BUFFERS, VERBOSE)
SELECT
	p.id, p.uuid, p.filename, p.title_cn, p.title_en, p.thumb_url, p.display_url,
	p.width, p.height, p.orientation, p.shot_time, p.aperture, p.shutter_speed,
	p.iso, p.like_count, p.view_count, p.download_count
FROM photos p
WHERE p.is_published = TRUE
  AND p.year = 2024
  AND p.month = 10
  AND p.orientation = 'landscape'
  AND (
    lower(COALESCE(p.title_cn, '')) LIKE '%天空%'
    OR lower(COALESCE(p.title_en, '')) LIKE '%sky%'
    OR lower(p.filename) LIKE '%sky%'
    OR EXISTS (
      SELECT 1
      FROM photo_tags pt
      JOIN tags t ON t.id = pt.tag_id
      WHERE pt.photo_id = p.id
        AND lower(t.name) LIKE '%天空%'
    )
  )
ORDER BY p.shot_time DESC, p.id DESC
LIMIT 30 OFFSET 0;

