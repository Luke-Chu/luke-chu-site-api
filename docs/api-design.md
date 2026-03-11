# 接口设计说明

## 基础信息

- 路由前缀：`/api/v1`
- 统一响应格式：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

错误时：

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

## 1) 健康检查

### 接口

- `GET /api/v1/health`

### 说明

返回服务状态，用于健康探活。

### curl 示例

```bash
curl -X GET "http://localhost:8080/api/v1/health"
```

---

## 2) 图片列表

### 接口

- `GET /api/v1/photos`

### Query 参数

- `q`：统一搜索词，支持空格、英文逗号、中文逗号、顿号分词
- `page`：页码，默认 `1`
- `pageSize`：每页数量，默认 `30`，最大 `60`
- `sort`：`shot_time|like_count|view_count|download_count|created_at`
- `order`：`asc|desc`
- `tags`：标签筛选
- `tagMode`：`any|all`
- `orientation`：`landscape|portrait|square`
- `year`
- `month`
- `category`

### curl 示例

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

## 3) 筛选项

### 接口

- `GET /api/v1/filters`

### 说明

用于列表页初始化筛选项，返回 `years/categories/orientations/tagTypes/tags`。

### curl 示例

```bash
curl -X GET "http://localhost:8080/api/v1/filters"
```

---

## 4) 图片详情

### 接口

- `GET /api/v1/photos/:uuid`

### Path 参数

- `uuid`：图片 UUID

### 说明

返回单张已发布图片详情及标签。详情字段包含 `exposureCompensation`（对应数据库 `exposure_compensation`）。

### curl 示例

```bash
curl -X GET "http://localhost:8080/api/v1/photos/550e8400-e29b-41d4-a716-446655440000"
```

### 错误语义

- `400`：UUID 非法
- `404`：图片不存在
- `500`：服务异常

---

## 5) 记录浏览行为

### 接口

- `POST /api/v1/photos/:uuid/view`

### 说明

带防刷窗口（10 分钟）：

- 同一 visitor 在窗口内重复请求不重复累加
- 返回 `counted` 标记本次是否计数

### curl 示例

```bash
curl -X POST "http://localhost:8080/api/v1/photos/550e8400-e29b-41d4-a716-446655440000/view"
```

### 错误语义

- `400`：UUID 非法
- `404`：图片不存在
- `500`：服务异常

---

## 6) 点赞

### 接口

- `POST /api/v1/photos/:uuid/like`

### 说明

- visitor hash 来自中间件
- 使用 `photo_likes(photo_id, visitor_hash)` 去重
- 首次点赞 `liked=true`，重复点赞 `liked=false`

### curl 示例

```bash
curl -X POST "http://localhost:8080/api/v1/photos/550e8400-e29b-41d4-a716-446655440000/like"
```

### 错误语义

- `400`：UUID 非法或 visitor 信息缺失
- `404`：图片不存在
- `500`：服务异常

---

## 7) 取消点赞

### 接口

- `POST /api/v1/photos/:uuid/unlike`

### 说明

- visitor hash 来自中间件
- 删除当前 visitor 点赞记录
- 若删除成功，`like_count` 回退（不小于 0）

### curl 示例

```bash
curl -X POST "http://localhost:8080/api/v1/photos/550e8400-e29b-41d4-a716-446655440000/unlike"
```

---

## 8) 记录下载行为

### 接口

- `POST /api/v1/photos/:uuid/download`

### 说明

带防刷窗口（30 分钟）：

- 同一 visitor 在窗口内重复请求不重复累加
- 返回最新 `downloadCount`、`downloadUrl` 和 `counted`

### curl 示例

```bash
curl -X POST "http://localhost:8080/api/v1/photos/550e8400-e29b-41d4-a716-446655440000/download"
```

### 错误语义

- `400`：UUID 非法
- `404`：图片不存在
- `500`：服务异常

---

更多示例请参考：[curl-examples.md](/docs/curl-examples.md)

