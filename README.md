# luke-chu-site-api

`luke-chu.com` 个人网站后端服务。当前阶段先支持摄影图库模块，后续逐步扩展其他业务模块。

## 技术栈

- Go
- Gin
- sqlx
- PostgreSQL
- Viper (yaml 配置)
- Zap (日志)
- go-playground/validator
- google/uuid

## 目录结构

```
luke-chu-site-api/
├── cmd/server/main.go
├── configs/
├── internal/
│   ├── app/
│   ├── config/
│   ├── db/
│   ├── model/
│   ├── repository/
│   ├── service/
│   ├── handler/
│   ├── dto/
│   ├── pkg/
│   └── constant/
├── api/openapi/
├── docs/
├── scripts/
└── test/
```

## 本地运行

1. 确认 PostgreSQL 可连接（默认读取 `configs/config.local.yaml`）。
2. 执行：

```bash
APP_ENV=local go run ./cmd/server
```

启动默认监听 `:8080`。

## 配置说明

- 配置文件位置：`configs/config.{env}.yaml`
- 通过环境变量 `APP_ENV` 选择配置文件：
  - `local` -> `configs/config.local.yaml`
  - `dev` -> `configs/config.dev.yaml`
  - `prod` -> `configs/config.prod.yaml`
- 未设置 `APP_ENV` 时默认 `local`

## 当前接口

- `GET /api/v1/health`
- `GET /api/v1/photos`
- `GET /api/v1/photos/:uuid`
- `POST /api/v1/photos/:uuid/view`
- `POST /api/v1/photos/:uuid/like`
- `POST /api/v1/photos/:uuid/download`
- `GET /api/v1/tags`

> 第一版图库 SQL 仍在演进，部分接口会返回 mock/降级数据。

## 构建

```bash
./scripts/build.sh
```

产物输出：`bin/luke-chu-site-api`

## 后续待办

- 完整图库表结构与迁移脚本
- 图片筛选/搜索 SQL 优化
- OpenAPI/Swagger 文档生成
- 集成测试完善

