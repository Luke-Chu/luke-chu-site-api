# 接口设计说明

## 基础信息

- 路由前缀：`/api/v1`
- 响应格式：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

失败示例：

```json
{
  "code": 40002,
  "message": "invalid uuid",
  "data": null
}
```

## 接口列表

- `GET /api/v1/health`
- `GET /api/v1/photos`
- `GET /api/v1/filters`
- `GET /api/v1/photos/:uuid`
- `POST /api/v1/photos/:uuid/view`
- `POST /api/v1/photos/:uuid/like`
- `POST /api/v1/photos/:uuid/unlike`
- `POST /api/v1/photos/:uuid/download`

---

## 1. 健康检查

- 方法与路径：`GET /api/v1/health`
- 说明：用于服务探活

```bash
curl -X GET "http://localhost:8080/api/v1/health"
```

---

## 2. 图片列表

- 方法与路径：`GET /api/v1/photos`
- 说明：支持关键词搜索、标签筛选、排序与分页

Query 参数：

- `q`：统一搜索词，支持空格/英文逗号/中文逗号/顿号分词
- `page`：页码，默认 `1`
- `pageSize`：每页数量，默认 `30`，最大 `60`
- `sort`：`shot_time|like_count|view_count|download_count|created_at`
- `order`：`asc|desc`
- `tags`：标签字符串，逗号分隔
- `tagMode`：`any|all`
- `orientation`：`landscape|portrait|square`
- `year`、`month`、`category`

性能基线说明：

- 当前列表查询已针对排序/筛选/关键词/标签关联补充索引迁移
- 可用 `docs/sql/photo-list-explain.sql` 做 `EXPLAIN ANALYZE` 验证执行计划

```bash
curl -G "http://localhost:8080/api/v1/photos" \
  --data-urlencode "q=天空,风筝" \
  --data-urlencode "page=2" \
  --data-urlencode "pageSize=20" \
  --data-urlencode "sort=view_count" \
  --data-urlencode "order=asc" \
  --data-urlencode "tags=风光,夕阳" \
  --data-urlencode "tagMode=all"
```

---

## 3. 筛选项

- 方法与路径：`GET /api/v1/filters`
- 说明：返回 years/categories/orientations/tagTypes/tags

```bash
curl -X GET "http://localhost:8080/api/v1/filters"
```

---

## 4. 图片详情

- 方法与路径：`GET /api/v1/photos/:uuid`
- Path 参数：`uuid`（图片 UUID）
- 说明：返回已发布图片详情和标签，包含 `exposureCompensation`

```bash
curl -X GET "http://localhost:8080/api/v1/photos/550e8400-e29b-41d4-a716-446655440000"
```

错误语义：

- `400`：UUID 非法
- `404`：图片不存在
- `500`：服务异常

---

## 5. 浏览行为

- 方法与路径：`POST /api/v1/photos/:uuid/view`
- 说明：10 分钟防刷窗口；窗口内重复请求不重复计数，返回 `counted=false`

```bash
curl -X POST "http://localhost:8080/api/v1/photos/550e8400-e29b-41d4-a716-446655440000/view"
```

---

## 6. 点赞

- 方法与路径：`POST /api/v1/photos/:uuid/like`
- 说明：visitor hash 来自中间件；同一 visitor 不重复点赞

```bash
curl -X POST "http://localhost:8080/api/v1/photos/550e8400-e29b-41d4-a716-446655440000/like"
```

---

## 7. 取消点赞

- 方法与路径：`POST /api/v1/photos/:uuid/unlike`
- 说明：删除当前 visitor 点赞记录，计数回退并保证不小于 0

```bash
curl -X POST "http://localhost:8080/api/v1/photos/550e8400-e29b-41d4-a716-446655440000/unlike"
```

---

## 8. 下载行为

- 方法与路径：`POST /api/v1/photos/:uuid/download`
- 说明：
  - 30 分钟防刷窗口；窗口内重复请求不重复计数，返回 `counted=false`
  - 返回 `downloadUrl` 为 OSS 预签名临时 URL（非裸公网 URL）

```bash
curl -X POST "http://localhost:8080/api/v1/photos/550e8400-e29b-41d4-a716-446655440000/download"
```

错误语义：

- `400`：UUID 非法
- `404`：图片不存在
- `500`：服务异常（含 OSS 预签名失败）

---

更多示例见：`docs/curl-examples.md`
