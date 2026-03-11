# curl 调用示例

基础地址：`http://localhost:8080`

## 1. 健康检查

```bash
curl -X GET "http://localhost:8080/api/v1/health"
```

## 2. 图片列表

### 2.1 基础查询

```bash
curl -X GET "http://localhost:8080/api/v1/photos"
```

### 2.2 关键词查询

```bash
curl -G "http://localhost:8080/api/v1/photos" \
  --data-urlencode "q=天空,风筝"
```

### 2.3 标签 + 排序 + 分页

```bash
curl -G "http://localhost:8080/api/v1/photos" \
  --data-urlencode "q=海边" \
  --data-urlencode "tags=风光,夕阳" \
  --data-urlencode "tagMode=all" \
  --data-urlencode "sort=view_count" \
  --data-urlencode "order=asc" \
  --data-urlencode "page=2" \
  --data-urlencode "pageSize=20"
```

## 3. 筛选项

```bash
curl -X GET "http://localhost:8080/api/v1/filters"
```

## 4. 图片详情

```bash
curl -X GET "http://localhost:8080/api/v1/photos/550e8400-e29b-41d4-a716-446655440000"
```

## 5. 记录浏览

```bash
curl -X POST "http://localhost:8080/api/v1/photos/550e8400-e29b-41d4-a716-446655440000/view"
```

## 6. 点赞

```bash
curl -X POST "http://localhost:8080/api/v1/photos/550e8400-e29b-41d4-a716-446655440000/like"
```

## 7. 取消点赞

```bash
curl -X POST "http://localhost:8080/api/v1/photos/550e8400-e29b-41d4-a716-446655440000/unlike"
```

## 8. 下载

```bash
curl -X POST "http://localhost:8080/api/v1/photos/550e8400-e29b-41d4-a716-446655440000/download"
```

说明：`download` 接口返回的 `downloadUrl` 是 OSS 预签名临时 URL，前端应直接使用该 URL 发起下载请求。
