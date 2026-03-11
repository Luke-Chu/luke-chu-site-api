# luke-chu-site-api 架构说明

## 分层设计

- `handler`: 参数绑定、调用 service、统一响应
- `service`: 业务编排与响应组装
- `repository`: SQL 查询与数据访问（sqlx + PostgreSQL）
- `db`: 数据库连接初始化与连接池
- `model`: 数据库实体映射
- `dto`: 请求/响应结构

调用链：`handler -> service -> repository -> db`

## 图片列表接口实现要点

- 过滤条件构建集中在 repository，复用到列表和总数查询。
- `q` 支持多关键词，采用 `EXISTS` 子查询匹配标签，避免主查询 join 造成重复行。
- 标签过滤支持 `tagMode=any|all`。
- 列表与标签分两步查询：
  1. 查询分页照片
  2. 按 photoIDs 批量查询标签并在 service 组装

## 筛选项接口实现要点

- `years/categories/orientations` 来自已发布照片聚合。
- `tags` 通过 `tags + photo_tags + photos` 关联查询，仅返回与已发布照片有关联的标签。
- `tagTypes` 固定返回 `subject|element|mood`。

