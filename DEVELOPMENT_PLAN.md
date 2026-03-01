# 微服务Demo项目发展计划

## 1. 项目结构设计

### 1.1 总体架构

项目将采用典型的微服务架构，使用RabbitMQ作为消息中间件，按照internal结构组织代码。

### 1.2 目录结构

```
eshop-rabbitmq-demo/
├── cmd/
│   ├── user-service/
│   │   └── main.go
│   ├── order-service/
│   │   └── main.go
│   ├── inventory-service/
│   │   └── main.go
│   ├── rbac-service/
│   │   └── main.go
│   └── log-service/
│       └── main.go
├── internal/
│   ├── user-service/
│   │   ├── api/
│   │   │   ├── handlers/
│   │   │   ├── middleware/
│   │   │   └── routes/
│   │   ├── domain/
│   │   │   ├── models/
│   │   │   └── repositories/
│   │   ├── service/
│   │   └── mq/
│   ├── order-service/
│   │   ├── api/
│   │   │   ├── handlers/
│   │   │   ├── middleware/
│   │   │   └── routes/
│   │   ├── domain/
│   │   │   ├── models/
│   │   │   └── repositories/
│   │   ├── service/
│   │   └── mq/
│   ├── inventory-service/
│   │   ├── api/
│   │   │   ├── handlers/
│   │   │   ├── middleware/
│   │   │   └── routes/
│   │   ├── domain/
│   │   │   ├── models/
│   │   │   └── repositories/
│   │   ├── service/
│   │   └── mq/
│   ├── rbac-service/
│   │   ├── api/
│   │   │   ├── handlers/
│   │   │   ├── middleware/
│   │   │   └── routes/
│   │   ├── domain/
│   │   │   ├── models/
│   │   │   └── repositories/
│   │   ├── service/
│   │   └── mq/
│   └── log-service/
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
│   ├── mq/
│   │   ├── event.go
│   │   └── mq.go
│   ├── config/
│   ├── database/
│   └── utils/
├── configs/
├── scripts/
├── Dockerfile
├── docker-compose.yml
├── go.work
└── README.md
```

## 2. 服务设计

### 2.1 User服务

**功能职责：**
- 用户注册、登录、注销
- 用户信息管理
- 用户个人资料维护
- 用户认证与授权

**API接口：**
- POST /api/users/register - 用户注册
- POST /api/users/login - 用户登录
- GET /api/users/profile - 获取用户资料
- PUT /api/users/profile - 更新用户资料
- POST /api/users/logout - 用户注销

**消息事件：**
- UserCreatedEvent - 用户创建事件
- UserUpdatedEvent - 用户更新事件
- UserDeletedEvent - 用户删除事件

### 2.2 Order服务

**功能职责：**
- 订单创建、查询、更新
- 订单状态管理
- 订单历史记录
- 订单与用户、库存的交互

**API接口：**
- POST /api/orders - 创建订单
- GET /api/orders - 获取订单列表
- GET /api/orders/{id} - 获取订单详情
- PUT /api/orders/{id} - 更新订单状态
- DELETE /api/orders/{id} - 取消订单

**消息事件：**
- OrderCreatedEvent - 订单创建事件
- OrderUpdatedEvent - 订单更新事件
- OrderCancelledEvent - 订单取消事件
- OrderCompletedEvent - 订单完成事件

### 2.3 库存服务

**功能职责：**
- 商品库存管理
- 库存变更记录
- 库存预警
- 与订单服务的交互

**API接口：**
- GET /api/inventory - 获取库存列表
- GET /api/inventory/{productId} - 获取商品库存
- PUT /api/inventory/{productId} - 更新商品库存
- POST /api/inventory - 添加新商品库存
- DELETE /api/inventory/{productId} - 删除商品库存

**消息事件：**
- InventoryUpdatedEvent - 库存更新事件
- InventoryLowEvent - 库存不足事件
- InventoryAddedEvent - 库存添加事件

### 2.4 RBAC服务

**功能职责：**
- 角色管理
- 权限管理
- 角色权限分配
- 基于角色的访问控制

**API接口：**
- POST /api/roles - 创建角色
- GET /api/roles - 获取角色列表
- PUT /api/roles/{id} - 更新角色
- DELETE /api/roles/{id} - 删除角色
- POST /api/permissions - 创建权限
- GET /api/permissions - 获取权限列表
- POST /api/roles/{id}/permissions - 为角色分配权限
- GET /api/users/{id}/roles - 获取用户角色
- POST /api/users/{id}/roles - 为用户分配角色

**消息事件：**
- RoleCreatedEvent - 角色创建事件
- RoleUpdatedEvent - 角色更新事件
- RoleDeletedEvent - 角色删除事件
- PermissionCreatedEvent - 权限创建事件
- PermissionUpdatedEvent - 权限更新事件
- PermissionDeletedEvent - 权限删除事件

### 2.5 Log服务

**功能职责：**
- 系统日志收集
- 日志分析与存储
- 日志查询与检索
- 异常监控

**API接口：**
- POST /api/logs - 提交日志
- GET /api/logs - 获取日志列表
- GET /api/logs/{id} - 获取日志详情
- GET /api/logs/query - 按条件查询日志

**消息事件：**
- LogEvent - 日志事件

## 3. 数据库设计

### 3.1 User服务数据库

**`users`表**
| 字段名 | 数据类型 | 约束 | 描述 |
|-------|---------|------|------|
| `id` | `UUID` | `PRIMARY KEY` | 用户ID |
| `username` | `VARCHAR(50)` | `UNIQUE NOT NULL` | 用户名 |
| `email` | `VARCHAR(100)` | `UNIQUE NOT NULL` | 邮箱 |
| `password_hash` | `VARCHAR(255)` | `NOT NULL` | 密码哈希 |
| `full_name` | `VARCHAR(100)` | | 全名 |
| `phone` | `VARCHAR(20)` | | 电话号码 |
| `created_at` | `TIMESTAMP` | `DEFAULT CURRENT_TIMESTAMP` | 创建时间 |
| `updated_at` | `TIMESTAMP` | `DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` | 更新时间 |

### 3.2 Order服务数据库

**`orders`表**
| 字段名 | 数据类型 | 约束 | 描述 |
|-------|---------|------|------|
| `id` | `UUID` | `PRIMARY KEY` | 订单ID |
| `customer_id` | `UUID` | `REFERENCES users(id)` | 客户ID |
| `total_amount` | `DECIMAL(10,2)` | `NOT NULL` | 总金额 |
| `status` | `VARCHAR(20)` | `NOT NULL` | 订单状态 |
| `created_at` | `TIMESTAMP` | `DEFAULT CURRENT_TIMESTAMP` | 创建时间 |
| `updated_at` | `TIMESTAMP` | `DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` | 更新时间 |

**`order_items`表**
| 字段名 | 数据类型 | 约束 | 描述 |
|-------|---------|------|------|
| `id` | `UUID` | `PRIMARY KEY` | 订单项ID |
| `order_id` | `UUID` | `REFERENCES orders(id)` | 订单ID |
| `product_id` | `UUID` | `NOT NULL` | 商品ID |
| `quantity` | `INT` | `NOT NULL` | 数量 |
| `unit_price` | `DECIMAL(10,2)` | `NOT NULL` | 单价 |

### 3.3 库存服务数据库

**`inventory`表**
| 字段名 | 数据类型 | 约束 | 描述 |
|-------|---------|------|------|
| `id` | `UUID` | `PRIMARY KEY` | 库存ID |
| `product_id` | `UUID` | `UNIQUE NOT NULL` | 商品ID |
| `quantity` | `INT` | `NOT NULL` | 库存数量 |
| `min_threshold` | `INT` | `NOT NULL` | 最小库存阈值 |
| `created_at` | `TIMESTAMP` | `DEFAULT CURRENT_TIMESTAMP` | 创建时间 |
| `updated_at` | `TIMESTAMP` | `DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` | 更新时间 |

**`inventory_history`表**
| 字段名 | 数据类型 | 约束 | 描述 |
|-------|---------|------|------|
| `id` | `UUID` | `PRIMARY KEY` | 历史记录ID |
| `product_id` | `UUID` | `NOT NULL` | 商品ID |
| `change_amount` | `INT` | `NOT NULL` | 变更数量 |
| `new_quantity` | `INT` | `NOT NULL` | 变更后数量 |
| `reason` | `VARCHAR(255)` | `NOT NULL` | 变更原因 |
| `created_at` | `TIMESTAMP` | `DEFAULT CURRENT_TIMESTAMP` | 创建时间 |

### 3.4 RBAC服务数据库

**`roles`表**
| 字段名 | 数据类型 | 约束 | 描述 |
|-------|---------|------|------|
| `id` | `UUID` | `PRIMARY KEY` | 角色ID |
| `name` | `VARCHAR(50)` | `UNIQUE NOT NULL` | 角色名称 |
| `description` | `VARCHAR(255)` | | 角色描述 |
| `created_at` | `TIMESTAMP` | `DEFAULT CURRENT_TIMESTAMP` | 创建时间 |
| `updated_at` | `TIMESTAMP` | `DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` | 更新时间 |

**`permissions`表**
| 字段名 | 数据类型 | 约束 | 描述 |
|-------|---------|------|------|
| `id` | `UUID` | `PRIMARY KEY` | 权限ID |
| `name` | `VARCHAR(50)` | `UNIQUE NOT NULL` | 权限名称 |
| `description` | `VARCHAR(255)` | | 权限描述 |
| `created_at` | `TIMESTAMP` | `DEFAULT CURRENT_TIMESTAMP` | 创建时间 |
| `updated_at` | `TIMESTAMP` | `DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` | 更新时间 |

**`role_permissions`表**
| 字段名 | 数据类型 | 约束 | 描述 |
|-------|---------|------|------|
| `role_id` | `UUID` | `REFERENCES roles(id)` | 角色ID |
| `permission_id` | `UUID` | `REFERENCES permissions(id)` | 权限ID |
| `PRIMARY KEY` | | `(role_id, permission_id)` | 复合主键 |

**`user_roles`表**
| 字段名 | 数据类型 | 约束 | 描述 |
|-------|---------|------|------|
| `user_id` | `UUID` | `REFERENCES users(id)` | 用户ID |
| `role_id` | `UUID` | `REFERENCES roles(id)` | 角色ID |
| `PRIMARY KEY` | | `(user_id, role_id)` | 复合主键 |

### 3.5 Log服务数据库

**`logs`表**
| 字段名 | 数据类型 | 约束 | 描述 |
|-------|---------|------|------|
| `id` | `UUID` | `PRIMARY KEY` | 日志ID |
| `service` | `VARCHAR(50)` | `NOT NULL` | 服务名称 |
| `level` | `VARCHAR(20)` | `NOT NULL` | 日志级别 |
| `message` | `TEXT` | `NOT NULL` | 日志消息 |
| `metadata` | `JSONB` | | 元数据 |
| `created_at` | `TIMESTAMP` | `DEFAULT CURRENT_TIMESTAMP` | 创建时间 |

## 4. 消息队列设计

### 4.1 事件定义

```go
// pkg/mq/event.go

package mq

// 基础事件结构
type Event struct {
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	Timestamp string          `json:"timestamp"`
	Source    string          `json:"source"`
}

// 用户相关事件
type UserCreatedEvent struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserUpdatedEvent struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserDeletedEvent struct {
	ID string `json:"id"`
}

// 订单相关事件
type OrderCreatedEvent struct {
	ID          string `json:"id"`
	CustomerID  string `json:"customer_id"`
	TotalAmount float64 `json:"total_amount"`
	Status      string `json:"status"`
}

type OrderUpdatedEvent struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type OrderCancelledEvent struct {
	ID string `json:"id"`
}

type OrderCompletedEvent struct {
	ID string `json:"id"`
}

// 库存相关事件
type InventoryUpdatedEvent struct {
	ProductID    string `json:"product_id"`
	ChangeAmount int    `json:"change_amount"`
	NewQuantity  int    `json:"new_quantity"`
}

type InventoryLowEvent struct {
	ProductID   string `json:"product_id"`
	CurrentQuantity int    `json:"current_quantity"`
	Threshold   int    `json:"threshold"`
}

type InventoryAddedEvent struct {
	ProductID   string `json:"product_id"`
	Quantity    int    `json:"quantity"`
}

// RBAC相关事件
type RoleCreatedEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoleUpdatedEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoleDeletedEvent struct {
	ID string `json:"id"`
}

type PermissionCreatedEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PermissionUpdatedEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PermissionDeletedEvent struct {
	ID string `json:"id"`
}

// 日志事件
type LogEvent struct {
	Service   string          `json:"service"`
	Level     string          `json:"level"`
	Message   string          `json:"message"`
	Metadata  json.RawMessage `json:"metadata"`
}
```

### 4.2 消息队列配置

- **Exchange名称**: `eshop-events`
- **Exchange类型**: `topic`
- **队列绑定规则**:
  - User服务: `user.*`
  - Order服务: `order.*`
  - 库存服务: `inventory.*`
  - RBAC服务: `rbac.*`
  - Log服务: `log.*`

## 5. 实现步骤

### 5.1 第一步：项目结构重构

1. 创建 `cmd` 目录，用于存放各服务的入口文件
2. 创建 `internal` 目录，按照服务类型组织代码
3. 创建 `pkg` 目录，存放共享代码
4. 移动现有代码到新结构中

### 5.2 第二步：实现User服务

1. 创建User服务目录结构
2. 实现用户注册、登录、注销功能
3. 实现用户信息管理功能
4. 实现用户认证与授权
5. 配置数据库连接
6. 实现消息队列交互

### 5.3 第三步：完善Order服务

1. 重构Order服务代码结构
2. 实现订单创建、查询、更新功能
3. 实现订单状态管理
4. 实现与User服务、库存服务的交互
5. 配置数据库连接
6. 完善消息队列交互

### 5.4 第四步：完善库存服务

1. 重构库存服务代码结构
2. 实现商品库存管理功能
3. 实现库存变更记录
4. 实现库存预警机制
5. 实现与Order服务的交互
6. 配置数据库连接
7. 完善消息队列交互

### 5.5 第五步：实现RBAC服务

1. 创建RBAC服务目录结构
2. 实现角色管理功能
3. 实现权限管理功能
4. 实现角色权限分配
5. 实现基于角色的访问控制
6. 配置数据库连接
7. 实现消息队列交互

### 5.6 第六步：完善Log服务

1. 重构Log服务代码结构
2. 实现系统日志收集
3. 实现日志分析与存储
4. 实现日志查询与检索
5. 实现异常监控
6. 配置数据库连接
7. 完善消息队列交互

### 5.7 第七步：集成测试

1. 编写单元测试
2. 编写集成测试
3. 进行系统测试
4. 性能测试

### 5.8 第八步：部署配置

1. 更新Dockerfile
2. 更新docker-compose.yml
3. 配置环境变量
4. 部署到测试环境
5. 部署到生产环境

## 6. 技术栈

- **语言**: Go 1.20+
- **Web框架**: Gin/Echo
- **数据库**: PostgreSQL
- **消息队列**: RabbitMQ
- **认证**: JWT
- **ORM**: GORM
- **配置管理**: Viper
- **日志**: Zap
- **容器化**: Docker
- **编排**: Docker Compose

## 7. 时间规划

| 阶段 | 任务 | 时间估计 |
|------|------|----------|
| 1 | 项目结构重构 | 1天 |
| 2 | 实现User服务 | 3天 |
| 3 | 完善Order服务 | 2天 |
| 4 | 完善库存服务 | 2天 |
| 5 | 实现RBAC服务 | 3天 |
| 6 | 完善Log服务 | 1天 |
| 7 | 集成测试 | 2天 |
| 8 | 部署配置 | 1天 |
| **总计** | | **15天** |

## 8. 风险评估

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 服务间依赖复杂 | 高 | 实现服务降级机制，确保单个服务故障不影响整体系统 |
| 消息队列延迟 | 中 | 实现消息重试机制，监控消息队列健康状态 |
| 数据库性能 | 中 | 实现数据库索引优化，考虑使用缓存 |
| 安全性 | 高 | 实现HTTPS，使用JWT认证，定期安全审计 |
| 扩展性 | 中 | 设计模块化架构，支持水平扩展 |

## 9. 监控与运维

- **健康检查**: 实现各服务的健康检查接口
- **日志监控**: 使用ELK stack收集和分析日志
- **性能监控**: 使用Prometheus和Grafana监控系统性能
- **告警机制**: 配置关键指标告警
- **CI/CD**: 实现自动化构建和部署流程

## 10. 总结

本发展计划详细描述了微服务Demo项目的结构设计、服务设计、数据库设计、消息队列设计以及实现步骤。通过按照internal结构组织项目，实现user服务、order服务、库存服务和rbac服务，将构建一个功能完整、结构清晰的微服务系统。该系统将具备良好的可扩展性、可维护性和可靠性，为后续的功能扩展和性能优化奠定基础。