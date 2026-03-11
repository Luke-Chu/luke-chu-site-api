# luke-chu-site-api 架构说明

## 分层结构

- `handler`: Gin HTTP 接口层，负责参数绑定、调用 service、输出统一响应。
- `service`: 业务编排层，承接业务流程与容错策略。
- `repository`: 数据访问层，使用 `sqlx` 访问 PostgreSQL。
- `db`: 数据库初始化与连接池管理。
- `model`: 数据库实体映射结构。
- `dto`: 请求与响应的数据结构。

## 调用链

`handler -> service -> repository -> db`

## 当前状态

- 已完成可运行骨架与路由注册。
- 图片与标签查询的真实 SQL 会在表结构落地后补充。
- 当前仓储层通过 `to_regclass` 检查表存在，不存在时返回 `ErrNotImplemented`，由 service 提供降级响应。

