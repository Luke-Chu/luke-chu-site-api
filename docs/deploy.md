# 服务器部署指南（Docker + GHCR）

本文档提供 `luke-chu-site-api` 在 Linux 服务器上的推荐部署方案：

- GitHub Actions 自动构建并推送镜像到 GHCR
- 服务器使用 `docker compose` 拉取镜像并启动服务

## 1. 前置条件

服务器需已安装：

- Docker Engine 24+
- Docker Compose Plugin 2+

快速检查：

```bash
docker --version
docker compose version
```

## 2. 准备部署文件

项目根目录已提供：

- `docker-compose.prod.yml`
- `.env.prod.example`

在服务器部署目录中创建：

```bash
mkdir -p /opt/luke-chu-site-api
cd /opt/luke-chu-site-api
```

把以下文件放到该目录（建议从仓库复制）：

- `docker-compose.prod.yml`
- `.env.prod`（由 `.env.prod.example` 改名并填真实值）

示例：

```bash
cp .env.prod.example .env.prod
chmod 600 .env.prod
```

注意：

- `.env.prod` 中严禁提交真实密钥到 Git 仓库。
- `OSS_ACCESS_KEY_SECRET`、数据库密码等应只存在于服务器本地。

## 3. GHCR 访问配置

若镜像仓库是私有的，需要先登录 GHCR：

```bash
echo "<GH_PAT>" | docker login ghcr.io -u <github-username> --password-stdin
```

说明：

- `<GH_PAT>` 需要包含 `read:packages` 权限。
- 若 GHCR 镜像已设置为 public，可跳过登录。

## 4. 首次部署

```bash
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d
docker compose -f docker-compose.prod.yml ps
```

查看日志：

```bash
docker compose -f docker-compose.prod.yml logs -f api
```

健康检查：

```bash
curl -X GET "http://127.0.0.1:8080/api/v1/health"
```

## 5. 更新发布

当 CI/CD 推送了新镜像后，服务器执行：

```bash
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d
```

## 6. 回滚

建议在 `docker-compose.prod.yml` 固定镜像 tag（例如 `v0.1.0`），需要回滚时改回旧 tag 后执行：

```bash
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d
```

当前模板支持在 `.env.prod` 中通过 `IMAGE_TAG` 控制版本，例如：

```env
IMAGE_TAG=v0.1.0
```

## 7. GitHub Actions 流程

- `.github/workflows/docker-image.yml`
  - `push main` / `push tag(v*)` / `workflow_dispatch` 触发
  - 构建多架构镜像（`linux/amd64`、`linux/arm64`）
  - 推送到 `ghcr.io/luke-chu/luke-chu-site-api`

## 8. 安全建议

- 不要在任何仓库文件中写死 AK/SK、数据库密码。
- 生产环境使用独立的 OSS RAM 用户，权限最小化。
- 定期轮换密钥。
- 建议在服务器侧限制 8080 端口来源（Nginx 反向代理 + 防火墙）。
