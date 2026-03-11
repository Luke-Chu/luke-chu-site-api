# luke-chu-site-api

`luke-chu.com` 个人网站后端服务。当前阶段重点支持摄影图库模块，后续持续扩展。

## 技术栈

- Go
- Gin
- sqlx
- PostgreSQL
- Viper (YAML 配置)
- Zap
- go-playground/validator
- google/uuid

## 本地运行

1. 确认 PostgreSQL 可连接（默认读取 `configs/config.local.yaml`）。
2. 启动服务：

```bash
APP_ENV=local go run ./cmd/server
```

默认监听 `:8080`。

## 配置说明

- 配置文件：`configs/config.{env}.yaml`
- 环境变量：`APP_ENV=local|dev|prod`
- 未设置 `APP_ENV` 时默认 `local`

本地默认数据库配置：

- host: `127.0.0.1`
- port: `5432`
- user: `admin`
- password: `1234`
- dbname: `luke-chu-site`

## 当前接口

- `GET /api/v1/health`
- `GET /api/v1/photos`
- `GET /api/v1/photos/:uuid`
- `POST /api/v1/photos/:uuid/view`
- `POST /api/v1/photos/:uuid/like`
- `POST /api/v1/photos/:uuid/download`
- `GET /api/v1/tags`
- `GET /api/v1/filters`

## 图片列表查询参数

`GET /api/v1/photos` 支持：

- `q`
- `page`
- `pageSize`
- `sort`
- `order`
- `tags`
- `tagMode` (`any` / `all`)
- `orientation`
- `year`
- `month`
- `category`

返回结构：

- `data.list`
- `data.pagination`
- `data.query`

## 筛选项接口

`GET /api/v1/filters` 返回：

- `years`
- `categories`
- `orientations`（含 `count`）
- `tagTypes`
- `tags`（按 `tag_type` 分组）

## 构建

```bash
./scripts/build.sh
```

构建产物：`bin/luke-chu-site-api`

## 后续计划

- 完整迁移脚本与索引优化
- 复杂搜索 SQL 和统计优化
- OpenAPI / Swagger 文档生成
- 集成测试完善

