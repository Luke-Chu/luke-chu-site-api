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
- `behavior_guard`：行为接口风控（IP 限流、可疑 UA 限流、风险日志）

## 下载签名链路

`POST /api/v1/photos/:uuid/download` 的处理流程：

1. Repository 处理下载计数与 30 分钟防刷窗口
2. Service 调用 OSS 预签名组件
3. 返回临时签名 `downloadUrl` 给前端

说明：AK/SK 不写入代码，SDK 从环境变量读取 `OSS_ACCESS_KEY_ID` / `OSS_ACCESS_KEY_SECRET`。

## 行为风控策略（第一版）

适用接口：

- `POST /api/v1/photos/:uuid/view`
- `POST /api/v1/photos/:uuid/like`
- `POST /api/v1/photos/:uuid/unlike`
- `POST /api/v1/photos/:uuid/download`

策略组合：

1. Repository 时间窗口防刷（view/download）
2. Middleware IP 固定窗口限流
3. Middleware 可疑 User-Agent 独立限流
4. 命中可疑规则或限流时输出 `warn` 风控日志（含 IP、UA、path、visitor_hash）

默认阈值（可配置）：

- `window_seconds=60`
- `ip_limit_per_window=120`
- `suspicious_ip_limit_per_window=20`

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

## CI/CD（GitHub Actions）

### CI

- 工作流：`.github/workflows/ci.yml`
- 触发：`push(main)`、`pull_request`
- 检查：`gofmt`、`go mod tidy` 一致性、`go vet`、`go test`、`go build`

### CD

- 工作流：`.github/workflows/release.yml`
- 触发：`push tag(v*)` 或手动触发
- 产物：多平台二进制（linux/darwin/windows）
- 发布：tag 触发时自动创建 GitHub Release 并上传资产包

### Docker 镜像发布

- 工作流：`.github/workflows/docker-image.yml`
- 镜像仓库：`ghcr.io/luke-chu/luke-chu-site-api`
- 构建平台：`linux/amd64`、`linux/arm64`
- 触发：`push main`、`push tag(v*)`、`workflow_dispatch`
