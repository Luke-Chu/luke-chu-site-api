# 架构说明

## 分层结构

- `handler`：参数绑定、基础校验、调用 service、返回统一响应
- `service`：业务编排与领域规则处理
- `repository`：使用 sqlx 进行手写 SQL 访问 PostgreSQL
- `db`：数据库连接与连接池初始化
- `dto`：请求/响应结构
- `model`：数据库实体映射

调用链：`handler -> service -> repository -> db`

## 行为接口策略

- 点赞/取消点赞：通过 `photo_likes` 去重并在事务内维护 `photos.like_count`
- 浏览计数：通过 `photo_views` 做 10 分钟窗口防刷
- 下载计数：通过 `photo_downloads` 做 30 分钟窗口防刷

## 测试策略

当前测试分两层：

1. 单元测试
- 覆盖关键词解析、排序归一化、分页处理、visitor hash
- 补充了请求参数归一化测试

2. 接口层集成测试
- 使用 `net/http/httptest` + Gin Router
- 通过 stub service 验证 handler/router 行为
- 不依赖真实 PostgreSQL，重点验证参数校验、状态码与响应结构

