# API 设计草案

## 路由前缀

- `/api/v1`

## 已实现接口

- `GET /api/v1/health`
  - 健康检查接口，返回服务状态。
- `GET /api/v1/photos`
  - 图片列表接口，支持分页和基础筛选参数骨架。
- `GET /api/v1/photos/:uuid`
  - 图片详情接口。
- `POST /api/v1/photos/:uuid/view`
  - 记录浏览行为。
- `POST /api/v1/photos/:uuid/like`
  - 记录点赞行为。
- `POST /api/v1/photos/:uuid/download`
  - 记录下载行为。
- `GET /api/v1/tags`
  - 标签列表接口。

## 响应约定

统一结构：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

错误时 `code != 0`，`data` 为 `null`。

