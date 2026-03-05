# 错误码文档

本文档定义了 eShop 微服务系统中所有 API 错误码及其含义。

## 错误码格式

错误码采用 4 位数字格式，按功能模块分类：

- **1xxx**: 通用业务错误
- **2xxx**: 权限相关错误

## 通用业务错误 (1xxx)

| 错误码 | 错误名称             | HTTP 状态码 | 描述             | 解决方案                     |
| ------ | -------------------- | ----------- | ---------------- | ---------------------------- |
| 1001   | `ErrProductNotFound` | 404         | 产品未找到       | 检查产品 ID 是否正确         |
| 1002   | `ErrStockNotEnough`  | 400         | 库存不足         | 减少订购数量或联系管理员     |
| 1003   | `ErrInvalidParams`   | 400         | 请求参数无效     | 检查请求参数格式和必填项     |
| 1004   | `ErrPaginationQuery` | 400         | 分页查询参数无效 | 检查 page 和 page_size 参数  |
| 1005   | `ErrUnauthorized`    | 401         | 未授权           | 需要登录或提供有效的访问令牌 |
| 1006   | `ErrUserNotFound`    | 404         | 用户未找到       | 检查用户 ID 是否正确         |
| 1007   | `ErrOrderNotFound`   | 404         | 订单未找到       | 检查订单 ID 是否正确         |
| 1008   | `ErrDuplicateOrder`  | 409         | 重复订单         | 避免重复提交相同订单         |
| 1009   | `ErrPaymentFailed`   | 400         | 支付失败         | 检查支付信息或稍后重试       |

## 认证相关错误 (101x-102x)

| 错误码 | 错误名称                       | HTTP 状态码 | 描述               | 解决方案                     |
| ------ | ------------------------------ | ----------- | ------------------ | ---------------------------- |
| 1010   | `ErrInvalidCredentials`        | 401         | 凭证无效           | 检查用户名和密码             |
| 1011   | `ErrEmailAlreadyRegistered`    | 409         | 邮箱已注册         | 使用其他邮箱或登录现有账户   |
| 1012   | `ErrUserAlreadyRegistered`     | 409         | 用户已注册         | 使用其他用户名或登录现有账户 |
| 1013   | `ErrNotFound`                  | 404         | 资源未找到         | 检查资源 ID 是否正确         |
| 1014   | `ErrAccountDisabled`           | 403         | 账户已禁用         | 联系管理员启用账户           |
| 1015   | `ErrWechatClientNotConfigured` | 500         | 微信客户端未配置   | 联系管理员配置微信登录       |
| 1016   | `ErrUsernameAlreadyExists`     | 409         | 用户名已存在       | 使用其他用户名               |
| 1017   | `ErrUnsupportedProvider`       | 400         | 不支持的认证提供商 | 检查认证提供商配置           |
| 1018   | `ErrIdentityAlreadyBound`      | 409         | 身份已绑定         | 该身份已绑定到其他账户       |
| 1019   | `ErrInvalidToken`              | 401         | 无效令牌           | 重新登录获取新令牌           |
| 1020   | `ErrTokenRevoked`              | 401         | 令牌已撤销         | 重新登录获取新令牌           |
| 1021   | `ErrGenerateAccessToken`       | 500         | 生成访问令牌失败   | 联系管理员检查系统配置       |
| 1022   | `ErrGenerateRefreshToken`      | 500         | 生成刷新令牌失败   | 联系管理员检查系统配置       |
| 1023   | `ErrSaveRefreshToken`          | 500         | 保存刷新令牌失败   | 联系管理员检查数据库连接     |
| 1024   | `ErrUnexpectedSigningMethod`   | 401         | 意外的签名方法     | 检查 JWT 配置                |
| 1025   | `ErrParseToken`                | 401         | 解析令牌失败       | 重新登录获取新令牌           |

## 权限相关错误 (2xxx)

| 错误码 | 错误名称                     | HTTP 状态码 | 描述             | 解决方案                   |
| ------ | ---------------------------- | ----------- | ---------------- | -------------------------- |
| 2001   | `ErrPermissionNotFound`      | 404         | 权限未找到       | 检查权限 ID 或名称是否正确 |
| 2002   | `ErrPermissionAlreadyExists` | 409         | 权限已存在       | 使用其他权限名称           |
| 2003   | `ErrInvalidRoleName`         | 400         | 无效的角色名称   | 检查角色名称格式           |
| 2004   | `ErrInsufficientPermissions` | 403         | 权限不足         | 联系管理员分配所需权限     |
| 2005   | `ErrCannotModifySystemRole`  | 403         | 不能修改系统角色 | 系统内置角色不能被修改     |
| 2006   | `ErrCannotDeleteSystemRole`  | 403         | 不能删除系统角色 | 系统内置角色不能被删除     |

## 错误响应格式

所有 API 错误响应遵循统一格式：

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

### 错误响应

```json
{
  "code": 1003,
  "message": "invalid parameters",
  "data": null
}
```

## HTTP 状态码映射

| HTTP 状态码 | 说明           | 常见错误码                         |
| ----------- | -------------- | ---------------------------------- |
| 200         | 请求成功       | 0                                  |
| 400         | 请求参数错误   | 1003, 1004, 2003                   |
| 401         | 未授权         | 1005, 1010, 1019, 1020, 1024, 1025 |
| 403         | 权限不足       | 1014, 2004, 2005, 2006             |
| 404         | 资源未找到     | 1001, 1006, 1007, 1013, 2001       |
| 409         | 资源冲突       | 1008, 1011, 1012, 1016, 1018, 2002 |
| 500         | 服务器内部错误 | 1009, 1015, 1021, 1022, 1023       |

## 错误处理最佳实践

### 客户端错误处理

1. **检查 HTTP 状态码**：首先根据 HTTP 状态码判断错误类型
2. **解析错误码**：根据响应中的 `code` 字段确定具体错误
3. **显示友好提示**：根据 `message` 字段向用户显示友好的错误提示
4. **提供解决方案**：参考本文档中的解决方案列，为用户提供解决建议

### 示例代码

```javascript
// JavaScript/TypeScript 示例
async function handleApiError(response) {
  const data = await response.json();

  switch (response.status) {
    case 400:
      console.error("请求参数错误:", data.message);
      break;
    case 401:
      console.error("未授权，请重新登录");
      // 跳转到登录页
      window.location.href = "/login";
      break;
    case 403:
      console.error("权限不足:", data.message);
      break;
    case 404:
      console.error("资源未找到:", data.message);
      break;
    case 500:
      console.error("服务器错误，请联系管理员");
      break;
    default:
      console.error("未知错误:", data.message);
  }

  // 根据错误码提供具体解决方案
  const solutions = {
    1003: "请检查请求参数格式",
    1005: "请先登录",
    1010: "用户名或密码错误",
    2004: "您没有权限执行此操作",
    2005: "系统内置角色不能被修改",
    2006: "系统内置角色不能被删除",
  };

  if (solutions[data.code]) {
    console.log("建议:", solutions[data.code]);
  }
}
```

```go
// Go 示例
func handleApiError(code int, message string) {
    switch code {
    case 1003:
        log.Println("请求参数错误:", message)
    case 1005, 1010, 1019, 1020, 1024, 1025:
        log.Println("认证失败:", message)
    case 2004:
        log.Println("权限不足:", message)
    default:
        log.Println("错误:", message)
    }
}
```

## 错误码扩展

如需添加新的错误码，请遵循以下规则：

1. **选择合适的分类**：根据功能选择 1xxx 或 2xxx 分类
2. **使用连续编号**：在同一分类内使用连续的编号
3. **提供清晰的描述**：错误名称和描述应该清晰易懂
4. **更新本文档**：添加新错误码时同步更新本文档
5. **提供解决方案**：为新错误码提供相应的解决方案

## 联系支持

如果遇到本文档未涵盖的错误，请联系技术支持团队：

- 邮箱：support@eshop.com
- 文档：https://docs.eshop.com/api/errors
