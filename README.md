# luke-chu-site-api

`luke-chu.com` 个人网站后端服务，当前以摄影图库模块为主，后续可继续扩展更多业务模块。

## 技术栈

- Go
- Gin
- sqlx
- PostgreSQL
- Viper（YAML 配置）
- Zap
- go-playground/validator
- google/uuid

## 本地运行

```bash
APP_ENV=local go run ./cmd/server
```

默认监听地址：`http://localhost:8080`

## 本地数据库默认配置

- host: `127.0.0.1`
- port: `5432`
- user: `admin`
- password: `1234`
- dbname: `luke-chu-site`

## 当前已实现接口

- `GET /api/v1/health`
- `GET /api/v1/photos`
- `GET /api/v1/filters`
- `GET /api/v1/photos/:uuid`
- `POST /api/v1/photos/:uuid/view`
- `POST /api/v1/photos/:uuid/like`
- `POST /api/v1/photos/:uuid/unlike`
- `POST /api/v1/photos/:uuid/download`

## 测试

### 运行全部测试

```bash
go test ./...
```

### 当前测试覆盖范围

- 工具层单元测试
  - `internal/pkg/search/keyword.go`
  - `internal/pkg/sort/sort.go`
  - `internal/pkg/pager/pager.go`
  - `internal/pkg/visitor/hash.go`
- 请求归一化单元测试
  - `internal/dto/request/photo_list_request.go`
- 接口层集成测试（基于 `httptest` + Gin Router + stub service）
  - `GET /api/v1/health`
  - `GET /api/v1/photos`
  - `GET /api/v1/filters`
  - `GET /api/v1/photos/:uuid`
  - `POST /api/v1/photos/:uuid/view`
  - `POST /api/v1/photos/:uuid/like`
  - `POST /api/v1/photos/:uuid/download`

## 文档

- 接口设计：`docs/api-design.md`
- curl 示例：`docs/curl-examples.md`
- 架构说明：`docs/architecture.md`
- OpenAPI 草稿：`api/openapi/openapi.yaml`

## 构建

```bash
./scripts/build.sh
```

构建产物：`bin/luke-chu-site-api`

