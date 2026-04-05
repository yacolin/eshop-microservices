# 客户端集成指南

本文档提供 eShop 微服务 API 的客户端集成示例和最佳实践。

## 目录

- [认证流程](#认证流程)
- [API 基础](#api-基础)
- [角色和权限管理](#角色和权限管理)
- [客户端示例](#客户端示例)
- [错误处理](#错误处理)
- [最佳实践](#最佳实践)

## 认证流程

### 1. 用户注册

```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure_password_123"
}
```

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "john_doe",
    "email": "john@example.com"
  }
}
```

### 2. 用户登录

```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "secure_password_123"
}
```

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600,
    "token_type": "Bearer"
  }
}
```

### 3. 刷新令牌

```bash
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 4. 用户登出

```bash
POST /api/v1/auth/logout
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## API 基础

### 请求头

所有需要认证的 API 请求都必须包含以下请求头：

```http
Authorization: Bearer {access_token}
Content-Type: application/json
```

### 响应格式

所有 API 响应遵循统一格式：

**成功响应：**

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

**错误响应：**

```json
{
  "code": 1003,
  "message": "invalid parameters",
  "data": null
}
```

## 角色和权限管理

### 获取角色列表

```bash
GET /api/v1/roles?page=1&page_size=20
Authorization: Bearer {access_token}
```

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "roles": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "name": "admin",
        "display_name": "管理员",
        "description": "系统超级管理员，拥有所有权限",
        "status": 1,
        "sort": 0,
        "is_system": true,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

### 创建角色（需要管理员权限）

```bash
POST /api/v1/roles
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "editor",
  "display_name": "编辑员",
  "description": "内容编辑员",
  "status": 1,
  "sort": 5,
  "is_system": false
}
```

**注意**：

- `is_system: true` 表示系统内置角色，通常由系统初始化时创建
- 系统内置角色（如 admin、customer、merchant、operator、system）不能被修改或删除
- 用户自定义角色应设置 `is_system: false`
- 只有超级管理员才能创建系统角色

### 为用户分配角色（需要管理员权限）

```bash
POST /api/v1/users/{user_id}/roles
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "role_id": "660e8400-e29b-41d4-a716-446655440001"
}
```

### 获取用户的角色列表

```bash
GET /api/v1/users/{user_id}/roles
Authorization: Bearer {access_token}
```

### 为角色分配权限（需要管理员权限）

```bash
POST /api/v1/roles/{role_id}/permissions
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "permission_ids": [
    "770e8400-e29b-41d4-a716-446655440002",
    "770e8400-e29b-41d4-a716-446655440003"
  ]
}
```

## 客户端示例

### JavaScript/TypeScript

```typescript
// API 客户端类
class EShopApiClient {
  private baseUrl: string;
  private accessToken: string | null = null;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  // 设置访问令牌
  setAccessToken(token: string) {
    this.accessToken = token;
  }

  // 通用请求方法
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (this.accessToken) {
      headers['Authorization'] = `Bearer ${this.accessToken}`;
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || '请求失败');
    }

    return data;
  }

  // 用户注册
  async register(username: string, email: string, password: string) {
    return this.request('/api/v1/auth/register', {
      method: 'POST',
      body: JSON.stringify({ username, email, password }),
    });
  }

  // 用户登录
  async login(username: string, password: string) {
    const response = await this.request<{
      code: number;
      message: string;
      data: {
        access_token: string;
        refresh_token: string;
        expires_in: number;
      };
    }>('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    });

    this.setAccessToken(response.data.access_token);
    return response.data;
  }

  // 获取角色列表
  async getRoles(page: number = 1, pageSize: number = 20) {
    return this.request(`/api/v1/roles?page=${page}&page_size=${pageSize}`);
  }

  // 创建角色
  async createRole(roleData: {
    name: string;
    display_name: string;
    description?: string;
    status?: number;
    sort?: number;
    is_system?: boolean;
  }) {
    return this.request('/api/v1/roles', {
      method: 'POST',
      body: JSON.stringify(roleData),
    });
  }

  // 为用户分配角色
  async assignRoleToUser(userId: string, roleId: string) {
    return this.request(`/api/v1/users/${userId}/roles`, {
      method: 'POST',
      body: JSON.stringify({ role_id: roleId }),
    });
  }
}

// 使用示例
const client = new EShopApiClient('https://api.eshop.com');

// 注册新用户
await client.register('john_doe', 'john@example.com', 'password123');

// 登录
const { access_token } = await client.login('john_doe', 'password123');

// 获取角色列表
const roles = await client.getRoles();
console.log('角色列表:', roles.data.roles);

// 创建新角色（需要管理员权限）
const newRole = await client.createRole({
  name: 'editor',
  display_name: '编辑员',
  description: '内容编辑员',
});
```

### Python

```python
import requests
from typing import Optional, Dict, Any

class EShopApiClient:
    def __init__(self, base_url: str):
        self.base_url = base_url
        self.access_token: Optional[str] = None

    def set_access_token(self, token: str):
        self.access_token = token

    def _get_headers(self) -> Dict[str, str]:
        headers = {'Content-Type': 'application/json'}
        if self.access_token:
            headers['Authorization'] = f'Bearer {self.access_token}'
        return headers

    def _request(self, endpoint: str, method: str = 'GET', data: Optional[Dict] = None) -> Dict[str, Any]:
        url = f'{self.base_url}{endpoint}'
        response = requests.request(
            method,
            url,
            json=data,
            headers=self._get_headers()
        )
        response.raise_for_status()
        return response.json()

    def register(self, username: str, email: str, password: str) -> Dict[str, Any]:
        return self._request('/api/v1/auth/register', 'POST', {
            'username': username,
            'email': email,
            'password': password
        })

    def login(self, username: str, password: str) -> Dict[str, Any]:
        response = self._request('/api/v1/auth/login', 'POST', {
            'username': username,
            'password': password
        })
        self.set_access_token(response['data']['access_token'])
        return response['data']

    def get_roles(self, page: int = 1, page_size: int = 20) -> Dict[str, Any]:
        return self._request(f'/api/v1/roles?page={page}&page_size={page_size}')

    def create_role(self, role_data: Dict[str, Any]) -> Dict[str, Any]:
        return self._request('/api/v1/roles', 'POST', role_data)

    def assign_role_to_user(self, user_id: str, role_id: str) -> Dict[str, Any]:
        return self._request(f'/api/v1/users/{user_id}/roles', 'POST', {
            'role_id': role_id
        })

# 使用示例
client = EShopApiClient('https://api.eshop.com')

# 注册新用户
client.register('john_doe', 'john@example.com', 'password123')

# 登录
login_data = client.login('john_doe', 'password123')

# 获取角色列表
roles = client.get_roles()
print('角色列表:', roles['data']['roles'])

# 创建新角色（需要管理员权限）
new_role = client.create_role({
    'name': 'editor',
    'display_name': '编辑员',
    'description': '内容编辑员',
})
```

### Go

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type EShopApiClient struct {
	baseURL    string
	accessToken string
	client     *http.Client
}

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func NewEShopApiClient(baseURL string) *EShopApiClient {
	return &EShopApiClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *EShopApiClient) SetAccessToken(token string) {
	c.accessToken = token
}

func (c *EShopApiClient) request(method, endpoint string, body interface{}) (*ApiResponse, error) {
	url := c.baseURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	return &apiResp, nil
}

func (c *EShopApiClient) Register(username, email, password string) (*ApiResponse, error) {
	return c.request("POST", "/api/v1/auth/register", map[string]string{
		"username": username,
		"email":    email,
		"password": password,
	})
}

func (c *EShopApiClient) Login(username, password string) (*LoginResponse, error) {
	resp, err := c.request("POST", "/api/v1/auth/login", map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, err
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(data, &loginResp); err != nil {
		return nil, err
	}

	c.SetAccessToken(loginResp.AccessToken)
	return &loginResp, nil
}

func (c *EShopApiClient) GetRoles(page, pageSize int) (*ApiResponse, error) {
	return c.request("GET", fmt.Sprintf("/api/v1/roles?page=%d&page_size=%d", page, pageSize), nil)
}

func main() {
	client := NewEShopApiClient("https://api.eshop.com")

	// 注册新用户
	_, err := client.Register("john_doe", "john@example.com", "password123")
	if err != nil {
		fmt.Println("注册失败:", err)
		return
	}

	// 登录
	loginData, err := client.Login("john_doe", "password123")
	if err != nil {
		fmt.Println("登录失败:", err)
		return
	}
	fmt.Println("登录成功，令牌:", loginData.AccessToken)

	// 获取角色列表
	roles, err := client.GetRoles(1, 20)
	if err != nil {
		fmt.Println("获取角色列表失败:", err)
		return
	}
	fmt.Println("角色列表:", roles.Data)
}
```

## 错误处理

### JavaScript/TypeScript 错误处理

```typescript
try {
  const response = await client.getRoles();
  console.log('角色列表:', response.data);
} catch (error) {
  if (error.message.includes('401')) {
    // 处理未授权错误
    console.error('未授权，请重新登录');
    // 跳转到登录页
    window.location.href = '/login';
  } else if (error.message.includes('403')) {
    // 处理权限不足错误
    console.error('权限不足:', error.message);
  } else {
    // 处理其他错误
    console.error('请求失败:', error.message);
  }
}
```

### Python 错误处理

```python
try:
    roles = client.get_roles()
    print('角色列表:', roles['data']['roles'])
except requests.exceptions.HTTPError as e:
    if e.response.status_code == 401:
        print('未授权，请重新登录')
    elif e.response.status_code == 403:
        print('权限不足')
    else:
        print(f'请求失败: {e}')
except Exception as e:
    print(f'发生错误: {e}')
```

## 最佳实践

### 1. 令牌管理

- **安全存储**：将访问令牌存储在安全的地方（如 HttpOnly Cookie 或内存中）
- **自动刷新**：在令牌过期前自动刷新
- **错误处理**：在收到 401 错误时自动重新登录

### 2. 错误重试

- **指数退避**：对于临时性错误，使用指数退避策略重试
- **最大重试次数**：设置最大重试次数避免无限重试
- **幂等性**：确保重试操作是幂等的

### 3. 性能优化

- **请求缓存**：对不常变化的数据进行缓存
- **批量请求**：将多个相关请求合并为批量请求
- **连接池**：使用连接池管理 HTTP 连接

### 4. 安全性

- **HTTPS**：始终使用 HTTPS 进行通信
- **CSRF 保护**：实现 CSRF 令牌保护
- **输入验证**：在客户端和服务器端都验证输入

## 联系支持

如果遇到集成问题，请联系技术支持：

- 邮箱：<support@eshop.com>
- 文档：<https://docs.eshop.com/api>
- 错误码文档：[ERROR\_CODES.md](ERROR_CODES.md)

