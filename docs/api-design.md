# API 设计说明

## 路由前缀

- `/api/v1`

## 已实现接口

- `GET /api/v1/health`
- `GET /api/v1/photos`
- `GET /api/v1/photos/:uuid`
- `POST /api/v1/photos/:uuid/view`
- `POST /api/v1/photos/:uuid/like`
- `POST /api/v1/photos/:uuid/download`
- `GET /api/v1/tags`
- `GET /api/v1/filters`

## 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

错误时 `code != 0`，`data = null`。

## GET /api/v1/photos

### Query 参数

- `q`: 统一搜索词，支持空格、`,`、`，`、`、` 分词
- `page`: 页码，默认 `1`
- `pageSize`: 每页数量，默认 `30`，最大 `60`
- `sort`: `shot_time|like_count|view_count|download_count|created_at`
- `order`: `asc|desc`
- `tags`: 标签字符串（分隔规则同 `q`）
- `tagMode`: `any|all`，默认 `any`
- `orientation`: `landscape|portrait|square`
- `year`
- `month`
- `category`

### 搜索与过滤逻辑

- 固定条件：`is_published = true`
- `q` 多关键词：组间 `AND`，组内（标题/文件名/标签）`OR`
- `tags + tagMode=any`：命中任意标签即可
- `tags + tagMode=all`：必须命中全部标签
- 排序字段与方向使用白名单归一化
- 分页使用 `LIMIT/OFFSET`
- 标签采用二次查询：先查照片，再按 photoIDs 批量查 tags

### 返回结构

- `data.list`: 照片列表（含 tags）
- `data.pagination`: 分页信息
- `data.query`: 归一化后的查询回显

## GET /api/v1/filters

用于列表页初始化筛选项，返回：

- `years`: 已发布图片年份（降序）
- `categories`: 已发布图片分类（去空）
- `orientations`: 三种方向及数量
- `tagTypes`: 固定 `subject|element|mood`
- `tags`: 按 `tag_type` 分组的标签列表（仅统计关联到已发布图片的标签）

