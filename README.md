# eshop-microservices

基于 DEVELOPMENT_PLAN.md 的微服务 Demo，技术栈：MySQL、Redis、GORM、Gin、amqp091-go、Viper。

## 目录结构

```
eshop-microservices/
├── cmd/
│   └── order-service/
│       └── main.go
├── internal/
│   └── order-service/
│       ├── api/
│       │   ├── handlers/
│       │   ├── middleware/
│       │   └── routes/
│       ├── domain/
│       │   ├── models/
│       │   └── repositories/
│       ├── service/
│       └── mq/
├── pkg/
│   ├── config/
│   ├── database/
│   ├── mq/
│   └── utils/
├── configs/
│   └── order-service.yaml
├── docker-compose.yml
├── go.mod
└── README.md
```

## 当前实现：order-service

- **API**
  - `POST /api/orders` 创建订单
  - `GET /api/orders` 订单列表（支持 `customer_id`、`page`、`page_size`）
  - `GET /api/orders/:id` 订单详情
  - `PUT /api/orders/:id` 更新订单状态
  - `DELETE /api/orders/:id` 取消订单
  - `GET /health` 健康检查

- **依赖**：MySQL（GORM）、Redis、RabbitMQ（可选，用于发布订单事件）

## 本地运行

### 方式一：使用本机 MySQL / Redis（不用 Docker）

1. 在本机安装并启动 **MySQL**、**Redis**。
2. 在 MySQL 中创建数据库：`CREATE DATABASE order_db;`
3. 复制并编辑本地配置（按需改密码）：
   ```bash
   # configs/order-service.local.yaml 中修改 mysql.password、redis.password
   ```
4. 指定本地配置并启动：
   ```powershell
   $env:CONFIG_PATH="configs/order-service.local.yaml"; go run ./cmd/order-service
   ```

### 方式二：使用 Docker 中的 MySQL / Redis

1. 启动依赖：
   ```bash
   docker-compose up -d
   ```
2. 运行服务（使用默认配置）：
   ```bash
   go run ./cmd/order-service
   ```

服务默认监听 `:8080`。可通过环境变量 `CONFIG_PATH` 指定配置文件路径。

## 多服务启动（单 go.mod 在根目录）

所有服务共用根目录的 `go.mod`，每个服务一个入口：`cmd/<服务名>/main.go`。

**本机启动指定服务：**
```powershell
go run ./cmd/order-service
go run ./cmd/user-service   # 新增服务后
```
指定配置：`$env:CONFIG_PATH="configs/order-service.local.yaml"; go run ./cmd/order-service`

**Docker 启动：** 使用同一 Dockerfile，通过构建参数 `SERVICE` 指定要编译的 cmd：
```powershell
docker-compose up -d --build
```
**修改代码后需重新构建镜像**（否则容器仍用旧镜像）：
```powershell
docker-compose up -d --build        # 推荐：会利用缓存，只改 .go 时代码层重建，不会重新 go mod download
```
Dockerfile 已按「先 COPY go.mod/go.sum 再 go mod download，最后 COPY 源码再 build」排序，**只改微服务代码时不会重新下载依赖**；只有改 go.mod/go.sum 时才会重新 download。启用 BuildKit 可进一步缓存模块目录：`$env:DOCKER_BUILDKIT=1; docker-compose up -d --build`。  
需要完全清缓存再构建时再用：`docker-compose build --no-cache order-service`
新增服务时：在 `docker-compose.yml` 里增加一个 service，设置 `build.args.SERVICE: user-service`、端口（如 8081）、`CONFIG_PATH` 和 `configs/user-service.docker.yaml` 即可。

## 配置说明

| 环境变量 | 说明 |
|----------|------|
| `CONFIG_PATH` | 配置文件路径，默认 `configs/order-service.yaml` |

配置文件支持 Viper 的 env 覆盖，例如：`MYSQL_HOST`、`SERVER_PORT` 等（点号用下划线替代）。
