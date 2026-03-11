# luke-chu-site-api

`luke-chu.com` 个人网站后端服务，当前以摄影图库模块为主，后续会扩展更多模块。

## 技术栈

- Go
- Gin
- sqlx
- PostgreSQL
- Viper（YAML 配置）
- Zap
- go-playground/validator
- google/uuid
- 阿里云 OSS Go SDK V2（下载预签名 URL）

## 本地运行

```bash
APP_ENV=local go run ./cmd/server
```

默认监听：`http://localhost:8080`

## 配置说明

### 1) APP_ENV

- `APP_ENV=local` -> `configs/config.local.yaml`
- `APP_ENV=dev` -> `configs/config.dev.yaml`
- `APP_ENV=prod` -> `configs/config.prod.yaml`

### 2) PostgreSQL（本地默认）

- host: `127.0.0.1`
- port: `5432`
- user: `admin`
- password: `1234`
- dbname: `luke-chu-site`

### 3) OSS 预签名下载配置

下载接口 `POST /api/v1/photos/:uuid/download` 现在返回 OSS 临时签名 URL。

- 代码中不会硬编码 AK/SK
- AK/SK 仅通过环境变量读取
- 推荐在本地 `.env` 或 CI/CD 密钥系统里注入

建议环境变量：

```bash
OSS_ACCESS_KEY_ID=your-access-key-id
OSS_ACCESS_KEY_SECRET=your-access-key-secret
OSS_SESSION_TOKEN=
OSS_BUCKET_NAME=luke-chu-site-photography
OSS_ENDPOINT=oss-cn-hongkong.aliyuncs.com
OSS_PUBLIC_BASE_URL=https://luke-chu-site-photography.oss-cn-hongkong.aliyuncs.com
OSS_REGION=cn-hongkong
OSS_PRESIGN_EXPIRE_SECONDS=300
```

说明：

- `OSS_ACCESS_KEY_ID` / `OSS_ACCESS_KEY_SECRET`：由 OSS SDK 环境凭证提供器读取
- `OSS_ENDPOINT`：外网访问域名（示例为香港地域）
- `OSS_PUBLIC_BASE_URL`：数据库 `original_url` 对应的公网基础地址，用于解析 object key
- `OSS_PRESIGN_EXPIRE_SECONDS`：预签名有效期，默认 300 秒
- 若 AK/SK 未配置，服务启动时会报错并退出，避免运行后才出现下载签名失败

## 当前已实现接口

- `GET /api/v1/health`
- `GET /api/v1/photos`
- `GET /api/v1/filters`
- `GET /api/v1/photos/:uuid`
- `POST /api/v1/photos/:uuid/view`
- `POST /api/v1/photos/:uuid/like`
- `POST /api/v1/photos/:uuid/unlike`
- `POST /api/v1/photos/:uuid/download`（返回临时签名 URL）

## 统一错误码（当前版本）

- `40000`：`invalid query params`
- `40002`：`invalid uuid`
- `40003`：`visitor hash missing`
- `40401`：`photo not found`
- `42901`：`too many behavior requests`
- `42902`：`suspicious behavior blocked`
- `50000`：`internal server error`
- `50001~50006`：各业务模块内部错误（列表/详情/点赞/下载/标签/筛选）

## 行为防刷策略（view/like/unlike/download）

已启用两层防刷：

1. 数据库时间窗口去重  
   `view=10分钟`，`download=30分钟`
2. 行为风控中间件  
   按 IP 固定窗口限流 + 可疑 User-Agent 单独限流 + 风险日志

默认配置（可在 `configs/*.yaml` 覆盖）：

```yaml
security:
  behavior:
    enabled: true
    window_seconds: 60
    ip_limit_per_window: 120
    suspicious_ip_limit_per_window: 20
```

可疑 UA 命中后会记录 `warn` 风控日志；若同窗口超过阈值会返回 `42902`。

## 查询性能基线（图库列表）

已新增索引迁移脚本：

- `internal/db/migrations/000001_add_photo_query_indexes.up.sql`
- `internal/db/migrations/000001_add_photo_query_indexes.down.sql`

核心覆盖点：

- 列表排序字段（`shot_time / like_count / view_count / download_count / created_at`）
- 常用筛选字段（`year / month / orientation / category`）
- 关键词搜索（`title_cn / title_en / filename / tags.name` 的 `ILIKE` 场景）
- 标签关联查询（`photo_tags(tag_id, photo_id)`）

性能验证模板：

- `docs/sql/photo-list-explain.sql`（`EXPLAIN ANALYZE` 示例）

## 测试

运行全部测试：

```bash
go test ./...
```

当前覆盖范围：

- 工具层单元测试：`keyword/sort/pager/visitor/oss`
- 请求归一化单元测试：`photo_list_request`
- 接口层集成测试（`httptest` + Gin Router + stub service）

## 文档

- 接口设计：`docs/api-design.md`
- curl 示例：`docs/curl-examples.md`
- 架构说明：`docs/architecture.md`
- OpenAPI 草稿：`api/openapi/openapi.yaml`

## 构建

```bash
./scripts/build.sh
```

产物路径：`bin/luke-chu-site-api`

## CI/CD（GitHub Actions）

已新增：

- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`

### CI 流程

触发条件：

- `push` 到 `main`
- 任意 `pull_request`

执行内容：

- `gofmt` 检查
- `go mod tidy` 一致性检查
- `go vet ./...`
- `go test ./...`
- `go build ./cmd/server`

### CD / 发布流程

触发条件：

- 推送 tag（`v*`，例如 `v0.1.0`）
- 手动触发 `workflow_dispatch`

执行内容：

- 构建多平台二进制：`linux/amd64`、`linux/arm64`、`darwin/amd64`、`darwin/arm64`、`windows/amd64`
- 自动打包产物（`tar.gz` / `zip`）
- tag 触发时自动创建 GitHub Release 并上传产物
