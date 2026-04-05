# eshop-microservices

一个完整的电商微服务系统 Demo，展示了现代微服务架构的最佳实践。

## 技术栈

- **语言**: Go 1.21+
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **缓存**: Redis 7
- **消息队列**: RabbitMQ 3.x
- **RPC**: gRPC + Protocol Buffers
- **配置管理**: Viper
- **日志**: Zap
- **容器化**: Docker + Docker Compose
- **API 文档**: Swagger/OpenAPI

## 核心特性

- ✅ **分布式事务**: 基于 Saga 模式的分布式事务协调，支持自动补偿
- ✅ **服务间通信**: gRPC 同步调用 + RabbitMQ 异步事件驱动
- ✅ **微服务架构**: 订单服务、库存服务、用户服务独立部署
- ✅ **API 规范**: 统一的 Swagger 文档和版本控制策略
- ✅ **认证授权**: JWT 认证 + Token 刷新机制
- ✅ **幂等性控制**: 基于 Redis 的请求幂等性保证
- ✅ **容器化部署**: Docker Compose 一键启动所有服务
- ✅ **可观测性**: 结构化日志 + Saga 状态追踪

## 项目亮点

### 1. Saga 分布式事务

实现了完整的 Saga 协调器，支持订单创建流程的分布式事务管理：

- 自动补偿机制：失败时自动回滚已执行的步骤
- 状态持久化：支持 Redis/Memory 两种存储方式
- 可观测性：提供 Saga 状态查询 API

### 2. 服务间通信

- **gRPC**: 订单服务同步调用库存服务（预占/扣减库存）
- **RabbitMQ**: 服务间异步事件通信，实现最终一致性

### 3. 完整的微服务

- **Order Service**: 订单管理、Saga 协调、事件发布
- **Inventory Service**: 产品管理、库存管理、gRPC 服务端
- **User Service**: 用户管理、JWT 认证

## 适用场景

本 Demo 适合用于：

- 📚 学习微服务架构设计
- 🎓 面试项目展示
- 💼 技术方案参考
- 🚀 微服务开发实践

## 目录结构

```
eshop-microservices/
├── api/
│   └── proto/
│       └── inventory.proto          # gRPC 协议定义
├── cmd/
│   ├── inventory-service/
│   │   ├── main.go
│   │   └── swagger.go
│   ├── order-service/
│   │   ├── main.go
│   │   └── swagger.go
│   └── user-service/
│       └── main.go
├── internal/
│   ├── inventory-service/
│   │   ├── api/
│   │   │   ├── dto/
│   │   │   ├── grpc/               # gRPC 服务端
│   │   │   ├── handlers/
│   │   │   └── routes/
│   │   ├── app/
│   │   ├── domain/
│   │   │   ├── models/
│   │   │   └── repositories/
│   │   ├── mq/
│   │   │   ├── consumer.go         # MQ 消费者
│   │   │   └── publisher.go        # MQ 发布者
│   │   └── service/
│   ├── order-service/
│   │   ├── api/
│   │   │   ├── dto/
│   │   │   ├── handlers/
│   │   │   └── routes/
│   │   ├── app/
│   │   ├── clients/
│   │   │   └── inventory_client.go # gRPC 客户端
│   │   ├── domain/
│   │   │   ├── models/
│   │   │   └── repositories/
│   │   ├── mq/
│   │   │   ├── consumer.go
│   │   │   ├── inventory_consumer.go
│   │   │   └── publisher.go
│   │   ├── saga/
│   │   │   └── create_order_saga.go # Saga 分布式事务
│   │   └── service/
│   └── user-service/
│       ├── api/
│       │   ├── dto/
│       │   ├── handlers/
│       │   │   ├── auth_handler.go
│       │   │   ├── permission_handler.go
│       │   │   ├── role_handler.go
│       │   │   └── user_handler.go
│       │   └── routes/
│       ├── app/
│       ├── domain/
│       │   ├── auth/
│       │   │   └── provider.go
│       │   ├── models/
│       │   │   ├── auth_token.go
│       │   │   ├── permission.go
│       │   │   ├── role.go
│       │   │   ├── user.go
│       │   │   └── user_identity.go
│       │   └── repositories/
│       │       ├── auth_token_repository.go
│       │       ├── permission_repository.go
│       │       ├── role_repository.go
│       │       ├── user_identity_repository.go
│       │       └── user_repository.go
│       ├── mq/
│       │   └── publisher.go
│       └── service/
│           ├── auth_service.go
│           ├── permission_service.go
│           ├── token_service.go
│           └── user_service.go
├── pkg/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── mysql.go
│   │   └── redis.go
│   ├── errcode/
│   │   └── errcode.go
│   ├── logger/
│   │   └── logger.go
│   ├── middleware/
│   │   ├── errorhandler.go
│   │   ├── idempotency.go
│   │   ├── jwtauth.go
│   │   ├── logger.go
│   │   ├── rbac.go
│   │   └── recovery.go
│   ├── mq/
│   │   ├── event.go
│   │   └── mq.go
│   ├── query/
│   │   └── query.go
│   ├── response/
│   │   └── response.go
│   ├── saga/
│   │   ├── memory_log.go
│   │   ├── redis_log.go
│   │   └── saga.go
│   └── utils/
│       ├── cryptopwd.go
│       ├── timestamp.go
│       └── utils.go
├── configs/
│   ├── inventory-service.yaml
│   ├── order-service.yaml
│   └── user-service.yaml
├── docs/
│   ├── CLIENT_INTEGRATION.md        # 客户端集成指南
│   ├── DEVELOPMENT_PLAN.md        # 开发计划
│   ├── ERROR_CODES.md              # 错误码文档
│   └── bugs.md                   # 已修复的 Bug 列表
├── nginx/
│   └── nginx.conf
├── scripts/
│   ├── install-deps.ps1
│   ├── mysql-init.sql
│   ├── permissions-init.sql        # 权限初始化脚本
│   ├── rbac-migration.sql         # RBAC 数据迁移脚本
│   └── user-roles-init.sql        # 用户角色初始化脚本
├── .dockerignore
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## 当前实现

### 1. Order Service（订单服务）

**端口**: 8080

**API**:

- `POST /api/v1/orders` 创建订单（使用 Saga 模式）
- `GET /api/v1/orders` 订单列表（支持分页、筛选）
- `GET /api/v1/orders/:id` 订单详情
- `PUT /api/v1/orders/:id` 更新订单状态
- `DELETE /api/v1/orders/:id` 取消订单
- `GET /api/v1/orders/saga/:saga_id` 获取 Saga 执行状态
- `GET /health` 健康检查
- `GET /swagger/*any` Swagger API 文档

**核心功能**:

- ✅ Saga 分布式事务（订单创建 + 库存预占）
- ✅ gRPC 客户端调用库存服务
- ✅ RabbitMQ 消息发布与消费
- ✅ 订单状态管理
- ✅ Swagger API 文档

**依赖**: MySQL、Redis、RabbitMQ、gRPC（库存服务）

---

### 2. Inventory Service（库存服务）

**端口**: 8081

**API**:

- `POST /api/v1/products` 创建产品
- `GET /api/v1/products` 产品列表
- `GET /api/v1/products/:id` 产品详情
- `PUT /api/v1/products/:id` 更新产品
- `DELETE /api/v1/products/:id` 删除产品
- `POST /api/v1/inventories` 创建库存
- `GET /api/v1/inventories` 库存列表
- `GET /api/v1/inventories/:id` 库存详情
- `PUT /api/v1/inventories/:id` 更新库存
- `DELETE /api/v1/inventories/:id` 删除库存
- `POST /api/v1/inventories/reserve` 预占库存
- `POST /api/v1/inventories/release` 释放库存
- `POST /api/v1/inventories/adjust` 调整库存
- `GET /api/v1/inventories/check-availability` 检查库存可用性
- `GET /api/v1/categories` 分类管理
- `GET /health` 健康检查
- `GET /swagger/*any` Swagger API 文档

**gRPC 服务**: 50051

- `ReserveStock` 预占库存
- `ConfirmDeduct` 确认扣减库存
- `ReleaseStock` 释放库存
- `CheckStockAvailability` 检查库存可用性

**核心功能**:

- ✅ 产品管理（CRUD）
- ✅ 库存管理（预占、释放、调整）
- ✅ 分类管理
- ✅ gRPC 服务端
- ✅ RabbitMQ 消费者（监听订单事件）
- ✅ Swagger API 文档

**依赖**: MySQL、Redis、RabbitMQ

---

### 3. User Service（用户服务）

**端口**: 8082

**API**:

- `POST /api/v1/users/register` 用户注册
- `POST /api/v1/users/login` 用户登录
- `GET /api/v1/users/profile` 获取用户资料
- `PUT /api/v1/users/profile` 更新用户资料
- `POST /api/v1/auth/refresh` 刷新 Token
- `GET /health` 健康检查
- `GET /api/v1/roles` 获取角色列表
- `GET /api/v1/roles/:id` 获取角色详情
- `POST /api/v1/roles` 创建角色（管理员）
- `PUT /api/v1/roles/:id` 更新角色（管理员）
- `DELETE /api/v1/roles/:id` 删除角色（管理员）
- `POST /api/v1/roles/:id/permissions` 为角色分配权限（管理员）
- `DELETE /api/v1/roles/:id/permissions` 移除角色的权限（管理员）
- `GET /api/v1/users/:user_id/roles` 获取用户的角色列表
- `POST /api/v1/users/:user_id/roles` 为用户分配角色（管理员）
- `DELETE /api/v1/users/:user_id/roles/:role_id` 移除用户的角色（管理员）
- `GET /api/v1/permissions` 获取权限列表
- `GET /api/v1/permissions/:id` 获取权限详情
- `POST /api/v1/permissions/check` 检查权限

**核心功能**:

- ✅ 用户注册与登录
- ✅ JWT 认证
- ✅ 密码加密
- ✅ Token 刷新机制
- ✅ RBAC 权限控制（基于角色的访问控制）
- ✅ 动态角色管理
- ✅ 细粒度权限控制（资源 + 操作）
  **依赖**: MySQL、Redis

---

## 技术架构

### 分布式事务

- **Saga 模式**: 订单创建流程的分布式事务协调
- **补偿机制**: 自动回滚失败的事务步骤
- **状态存储**: Redis/Memory 两种存储方式
- **可观测性**: 提供 Saga 状态查询 API

### 服务间通信

- **gRPC**: 订单服务 → 库存服务（同步调用）
- **RabbitMQ**: 服务间异步事件通信
  - 订单服务发布: `order.created`, `order.cancelled`, `order.confirmed`
  - 库存服务监听: 预占/释放库存
  - 支付服务（待实现）: `payment.completed`, `payment.failed`

### 中间件

- **JWT 认证**: 统一的用户认证中间件
- **RBAC 权限控制**: 基于角色的访问控制，支持动态角色管理
- **RBAC 中间件**: RequireRole、RequireAdmin、RequireMerchant（位于 pkg/middleware）
- **幂等性**: 基于 Redis 的请求幂等性控制
- **错误处理**: 统一的错误码和响应格式
- **日志**: Zap 结构化日志
- **恢复**: Panic 恢复中间件

## 本地运行

### 前置要求

- Go 1.21+
- MySQL 8.0+
- Redis 7+
- RabbitMQ 3.x
- Protocol Buffers compiler (protoc)

### 方式一：使用 Docker Compose（推荐）

1. 启动所有服务：

   ```bash
   docker-compose up -d
   ```

2. 查看服务状态：

   ```bash
   docker-compose ps
   ```

3. 查看日志：
   ```bash
   docker-compose logs -f
   ```

**服务访问**:

- Order Service: http://localhost:8080
- Inventory Service: http://localhost:8081
- User Service: http://localhost:8082
- Inventory gRPC: localhost:50051
- RabbitMQ Management: http://localhost:15673 (guest/guest)
- Swagger 文档:
  - Order Service: http://localhost:8080/swagger/index.html
  - Inventory Service: http://localhost:8081/swagger/index.html

---

### 方式二：本地运行服务（开发模式）

#### 1. 启动依赖服务

```bash
# 仅启动 MySQL、Redis、RabbitMQ
docker-compose up -d mysql redis rabbitmq
```

#### 2. 初始化数据库

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE order_db; CREATE DATABASE inventory_db; CREATE DATABASE user_db;"

# 执行初始化脚本
mysql -u root -p order_db < scripts/mysql-init.sql
```

#### 3. 运行服务

**Order Service**:

```powershell
$env:CONFIG_PATH="configs/order-service.yaml"
go run ./cmd/order-service
```

**Inventory Service**:

```powershell
$env:CONFIG_PATH="configs/inventory-service.yaml"
go run ./cmd/inventory-service
```

**User Service**:

```powershell
$env:CONFIG_PATH="configs/user-service.yaml"
go run ./cmd/user-service
```

#### 4. 生成 gRPC 代码（如需修改 proto）

```bash
# 安装 protoc 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 生成代码
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       api/proto/inventory.proto
```

---

### 方式三：使用本机 MySQL / Redis（不用 Docker）

1. 在本机安装并启动 **MySQL**、**Redis**、**RabbitMQ**。
2. 在 MySQL 中创建数据库：
   ```bash
   CREATE DATABASE order_db;
   CREATE DATABASE inventory_db;
   CREATE DATABASE user_db;
   ```
3. 复制并编辑本地配置（按需改密码）：
   ```bash
   # configs/order-service.local.yaml 中修改 mysql.password、redis.password
   ```
4. 指定本地配置并启动：
   ```powershell
   $env:CONFIG_PATH="configs/order-service.local.yaml"; go run ./cmd/order-service
   ```

---

## 多服务启动（单 go.mod 在根目录）

所有服务共用根目录的 `go.mod`，每个服务一个入口：`cmd/<服务名>/main.go`。

**本机启动指定服务：**

```powershell
go run ./cmd/order-service
go run ./cmd/inventory-service
go run ./cmd/user-service
```

指定配置：`$env:CONFIG_PATH="configs/order-service.yaml"; go run ./cmd/order-service`

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

## 配置说明

### 环境变量

| 环境变量                 | 说明                             | 默认值                   |
| ------------------------ | -------------------------------- | ------------------------ |
| `CONFIG_PATH`            | 配置文件路径                     | `configs/<service>.yaml` |
| `INVENTORY_SERVICE_ADDR` | 库存服务 gRPC 地址（订单服务用） | `localhost:50051`        |

### 配置文件

每个服务都有独立的配置文件：

- `configs/order-service.yaml` - 订单服务配置
- `configs/inventory-service.yaml` - 库存服务配置
- `configs/user-service.yaml` - 用户服务配置

配置文件支持 Viper 的 env 覆盖，例如：`MYSQL_HOST`、`SERVER_PORT` 等（点号用下划线替代）。

### 配置示例

**Order Service**:

```yaml
server:
  port: 8080
  mode: debug

mysql:
  host: localhost
  port: 3306
  user: root
  password: root
  dbname: order_db

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

rabbitmq:
  url: amqp://guest:guest@localhost:10672/
  exchange: eshop-events

jwt:
  secret: "your-secret-key"
```

**Inventory Service**:

```yaml
server:
  port: 8081
  mode: debug

grpc:
  port: 50051

mysql:
  host: localhost
  port: 3306
  user: root
  password: root
  dbname: inventory_db

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

rabbitmq:
  url: amqp://guest:guest@localhost:10672/
  exchange: eshop-events
```

## Docker Compose 常用操作

### 启动服务

```bash
# 启动所有服务（后台运行）
docker-compose up -d

# 重启所有服务
docker-compose restart

# 重启某个服务
docker-compose restart <service_name>


# 启动指定服务
docker-compose up -d order-service

# 重启指定服务
docker-compose restart order-service

# 启动并重新构建镜像
docker-compose up -d --build
```

### 停止服务

```bash
# 停止所有服务
docker-compose down

# 停止指定服务
docker-compose stop order-service
```

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs

# 查看指定服务日志
docker-compose logs order-service

# 持续跟踪日志（实时输出）
docker-compose logs -f order-service

# 查看最近的 N 行日志
docker-compose logs --tail=100 order-service

# 禁用颜色输出
docker-compose logs --no-color order-service
```

### 查看服务状态

```bash
# 查看所有服务状态
docker-compose ps

# 查看指定服务状态
docker-compose ps order-service
```

### 进入容器

```bash
# 进入指定服务容器
docker-compose exec order-service sh

# 进入容器并执行命令
docker-compose exec order-service ls -la
```

### 其他操作

```bash
# 构建镜像（不启动）
docker-compose build order-service

# 构建镜像（不使用缓存）
docker-compose build --no-cache order-service

# 查看配置
docker-compose config

# 查看容器资源使用情况
docker stats $(docker-compose ps -q)
```

---

## 项目特性

### 1. 分布式事务（Saga 模式）

订单创建流程使用 Saga 模式保证分布式事务一致性：

```
1. 创建订单 → 失败 → 删除订单
2. 预占库存 → 失败 → 释放库存
3. 确认订单 → 失败 → 恢复订单状态
```

**查询 Saga 状态**:

```bash
curl http://localhost:8080/api/v1/orders/saga/{saga_id}
```

### 2. 服务间通信

**gRPC 同步调用**:

- 订单服务 → 库存服务（预占/扣减库存）

**RabbitMQ 异步事件**:

- 订单服务发布: `order.created`, `order.cancelled`, `order.confirmed`
- 库存服务监听并处理库存变更

### 3. API 文档

所有服务都提供 Swagger 文档：

- Order Service: http://localhost:8080/swagger/index.html
- Inventory Service: http://localhost:8081/swagger/index.html

### 4. 中间件

- **JWT 认证**: 统一的用户认证
- **幂等性**: 基于 Redis 的请求幂等性控制
- **错误处理**: 统一的错误码和响应格式
- **日志**: Zap 结构化日志
- **恢复**: Panic 恢复中间件

---

## 使用示例

### 创建订单（Saga 模式）

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "123e4567-e89b-12d3-a456-426614174000",
    "currency": "CNY",
    "items": [
      {
        "product_id": "product-1",
        "quantity": 2,
        "unit_price": 10000
      }
    ]
  }'
```

**响应**:

```json
{
  "id": "order-id",
  "customer_id": "123e4567-e89b-12d3-a456-426614174000",
  "total_amount": 20000,
  "currency": "CNY",
  "status": "confirmed",
  "items": [...]
}
```

### 查询 Saga 状态

```bash
curl http://localhost:8080/api/v1/orders/saga/{saga_id}
```

**响应**:

```json
{
  "id": "saga-id",
  "name": "create_order",
  "status": "succeeded",
  "steps": [
    {
      "name": "create_order",
      "status": "succeeded"
    },
    {
      "name": "reserve_stock",
      "status": "succeeded"
    },
    {
      "name": "confirm_order",
      "status": "succeeded"
    }
  ]
}
```

### 创建产品

```bash
curl -X POST http://localhost:8081/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15",
    "description": "Apple iPhone 15",
    "price": 599900,
    "sku": "IPHONE15"
  }'
```

### 创建库存

```bash
curl -X POST http://localhost:8081/api/v1/inventories \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "product-id",
    "quantity": 100,
    "threshold": 10
  }'
```

### 用户注册

```bash
curl -X POST http://localhost:8082/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 用户登录

```bash
curl -X POST http://localhost:8082/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

---

## 开发指南

### 添加新服务

1. 创建服务目录结构：

   ```
   internal/<service-name>/
   ├── api/
   ├── app/
   ├── domain/
   ├── mq/
   └── service/
   ```

2. 创建入口文件：

   ```
   cmd/<service-name>/
   ├── main.go
   └── swagger.go
   ```

3. 创建配置文件：

   ```
   configs/<service-name>.yaml
   ```

4. 在 `docker-compose.yml` 中添加服务配置

### 添加新 API

1. 在 `api/dto/` 中定义 DTO
2. 在 `service/` 中实现业务逻辑
3. 在 `api/handlers/` 中创建处理器
4. 在 `api/routes/` 中注册路由
5. 添加 Swagger 注释

### 测试

```bash
# 运行所有测试
go test ./...

# 运行指定服务的测试
go test ./internal/order-service/...

# 运行测试并显示覆盖率
go test -cover ./...
```

---

## 故障排查

### 服务无法启动

1. 检查配置文件路径：

   ```bash
   echo $CONFIG_PATH
   ```

2. 检查数据库连接：

   ```bash
   mysql -u root -p -h localhost
   ```

3. 检查 RabbitMQ 连接：
   ```bash
   # 访问管理界面
   http://localhost:15673
   ```

### Saga 执行失败

1. 查询 Saga 状态：

   ```bash
   curl http://localhost:8080/api/v1/orders/saga/{saga_id}
   ```

2. 查看服务日志：

   ```bash
   docker-compose logs order-service
   ```

3. 检查库存服务：
   ```bash
   curl http://localhost:8081/health
   ```

### 消息队列问题

1. 查看 RabbitMQ 管理界面：

   ```
   http://localhost:15673
   ```

2. 检查队列状态：
   - 连接状态
   - 队列绑定
   - 消息积压

3. 查看消费者日志：
   ```bash
   docker-compose logs inventory-service
   ```

---

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

---

## 许可证

本项目采用 MIT 许可证 - 查看 LICENSE 文件了解详情

---

## 联系方式

- 项目链接: [https://github.com/yourusername/eshop-microservices](https://github.com/yourusername/eshop-microservices)
