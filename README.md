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
