# 架构说明

## 分层结构

项目采用典型三层调用链：

- Handler：参数绑定、输入校验、统一响应
- Service：业务编排与错误语义转换
- Repository：基于 sqlx 的手写 SQL 访问 PostgreSQL

请求路径：`HTTP -> Gin Router -> Middleware -> Handler -> Service -> Repository -> PostgreSQL`

## 中间件

- `cors`：本地开发跨域支持
- `logger`：请求方法、路径、状态码、耗时、IP
- `recovery`：panic 保护并返回统一错误 JSON
- `visitor`：基于 IP + UA + Accept-Language 生成 `visitor_hash`

## 下载签名链路

`POST /api/v1/photos/:uuid/download` 的处理流程：

1. Repository 处理下载计数与 30 分钟防刷窗口
2. Service 调用 OSS 预签名组件
3. 返回临时签名 `downloadUrl` 给前端

说明：AK/SK 不写入代码，SDK 从环境变量读取 `OSS_ACCESS_KEY_ID` / `OSS_ACCESS_KEY_SECRET`。

## 查询性能策略（第一版）

围绕 `GET /api/v1/photos` 与 `GET /api/v1/filters`，当前已落地：

1. 排序字段索引：`shot_time/like_count/view_count/download_count/created_at`（发布态部分索引）
2. 过滤字段索引：`year/month/orientation/category`（发布态部分索引）
3. 关键词检索索引：`title_cn/title_en/filename/tags.name` 的 `lower(...) + pg_trgm` 索引
4. 标签关联索引：`photo_tags(tag_id, photo_id)`

相关文件：

- `internal/db/migrations/000001_add_photo_query_indexes.up.sql`
- `internal/db/migrations/000001_add_photo_query_indexes.down.sql`
- `docs/sql/photo-list-explain.sql`

## 测试策略

- 单元测试：聚焦工具函数和请求归一化逻辑
- 接口集成测试：`httptest` + Gin Router + stub service，不依赖真实 PostgreSQL
