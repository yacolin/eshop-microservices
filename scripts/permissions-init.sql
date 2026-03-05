-- ====================================
-- 权限初始化脚本
-- ====================================
-- 用途: 初始化系统默认权限、角色和角色权限映射
-- 使用方式: mysql -u root -p user_db < scripts/permissions-init.sql
-- ====================================

USE user_db;

-- 清空现有权限和角色权限数据（可选，根据需要注释掉）
-- TRUNCATE TABLE role_permissions;
-- TRUNCATE TABLE roles;
-- DELETE FROM permissions WHERE id IS NOT NULL;

-- ====================================
-- 插入默认角色
-- ====================================
INSERT INTO roles (id, name, display_name, description, status, sort, is_system, created_at, updated_at) VALUES
(UUID(), 'admin', '管理员', '系统超级管理员，拥有所有权限', 1, 0, 1, NOW(), NOW()),
(UUID(), 'customer', '普通用户', '普通用户，基础业务权限', 1, 1, 1, NOW(), NOW()),
(UUID(), 'merchant', '商家', '商家用户，管理商品和订单', 1, 2, 1, NOW(), NOW()),
(UUID(), 'operator', '运营人员', '运营人员，审核和管理业务', 1, 3, 1, NOW(), NOW()),
(UUID(), 'system', '系统用户', '系统内部用户，用于系统间调用', 1, 4, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE
    display_name = VALUES(display_name),
    description = VALUES(description),
    updated_at = NOW();

-- ====================================
-- 插入默认权限
-- ====================================
INSERT INTO permissions (id, name, display_name, description, resource, action, category, sort, status, created_at, updated_at) VALUES
-- 用户相关权限 (admin 分类)
(UUID(), 'user:create', '创建用户', '创建新用户', 'user', 'create', 'admin', 0, 1, NOW(), NOW()),
(UUID(), 'user:read', '查看用户', '查看用户详情', 'user', 'read', 'admin', 1, 1, NOW(), NOW()),
(UUID(), 'user:update', '更新用户', '更新用户信息', 'user', 'update', 'admin', 2, 1, NOW(), NOW()),
(UUID(), 'user:delete', '删除用户', '删除用户', 'user', 'delete', 'admin', 3, 1, NOW(), NOW()),
(UUID(), 'user:list', '用户列表', '查看用户列表', 'user', 'list', 'admin', 4, 1, NOW(), NOW()),

-- 订单相关权限 (business 分类)
(UUID(), 'order:create', '创建订单', '创建新订单', 'order', 'create', 'business', 0, 1, NOW(), NOW()),
(UUID(), 'order:read', '查看订单', '查看订单详情', 'order', 'read', 'business', 1, 1, NOW(), NOW()),
(UUID(), 'order:update', '更新订单', '更新订单状态', 'order', 'update', 'business', 2, 1, NOW(), NOW()),
(UUID(), 'order:delete', '删除订单', '删除订单', 'order', 'delete', 'business', 3, 1, NOW(), NOW()),
(UUID(), 'order:list', '订单列表', '查看订单列表', 'order', 'list', 'business', 4, 1, NOW(), NOW()),
(UUID(), 'order:approve', '审批订单', '审批订单', 'order', 'approve', 'admin', 5, 1, NOW(), NOW()),
(UUID(), 'order:reject', '拒绝订单', '拒绝订单', 'order', 'reject', 'admin', 6, 1, NOW(), NOW()),

-- 产品相关权限 (business 分类)
(UUID(), 'product:create', '创建产品', '创建新产品', 'product', 'create', 'business', 0, 1, NOW(), NOW()),
(UUID(), 'product:read', '查看产品', '查看产品详情', 'product', 'read', 'business', 1, 1, NOW(), NOW()),
(UUID(), 'product:update', '更新产品', '更新产品信息', 'product', 'update', 'business', 2, 1, NOW(), NOW()),
(UUID(), 'product:delete', '删除产品', '删除产品', 'product', 'delete', 'business', 3, 1, NOW(), NOW()),
(UUID(), 'product:list', '产品列表', '查看产品列表', 'product', 'list', 'business', 4, 1, NOW(), NOW()),

-- 库存相关权限 (business 分类)
(UUID(), 'inventory:create', '创建库存', '创建库存记录', 'inventory', 'create', 'business', 0, 1, NOW(), NOW()),
(UUID(), 'inventory:read', '查看库存', '查看库存详情', 'inventory', 'read', 'business', 1, 1, NOW(), NOW()),
(UUID(), 'inventory:update', '更新库存', '更新库存数量', 'inventory', 'update', 'business', 2, 1, NOW(), NOW()),
(UUID(), 'inventory:delete', '删除库存', '删除库存记录', 'inventory', 'delete', 'business', 3, 1, NOW(), NOW()),
(UUID(), 'inventory:list', '库存列表', '查看库存列表', 'inventory', 'list', 'business', 4, 1, NOW(), NOW()),

-- 角色权限管理 (system 分类)
(UUID(), 'role:create', '创建角色', '创建新角色', 'role', 'create', 'system', 0, 1, NOW(), NOW()),
(UUID(), 'role:read', '查看角色', '查看角色详情', 'role', 'read', 'system', 1, 1, NOW(), NOW()),
(UUID(), 'role:update', '更新角色', '更新角色信息', 'role', 'update', 'system', 2, 1, NOW(), NOW()),
(UUID(), 'role:delete', '删除角色', '删除角色', 'role', 'delete', 'system', 3, 1, NOW(), NOW()),
(UUID(), 'role:list', '角色列表', '查看角色列表', 'role', 'list', 'system', 4, 1, NOW(), NOW()),

-- 权限管理 (system 分类)
(UUID(), 'permission:create', '创建权限', '创建新权限', 'permission', 'create', 'system', 0, 1, NOW(), NOW()),
(UUID(), 'permission:read', '查看权限', '查看权限详情', 'permission', 'read', 'system', 1, 1, NOW(), NOW()),
(UUID(), 'permission:update', '更新权限', '更新权限信息', 'permission', 'update', 'system', 2, 1, NOW(), NOW()),
(UUID(), 'permission:delete', '删除权限', '删除权限', 'permission', 'delete', 'system', 3, 1, NOW(), NOW()),
(UUID(), 'permission:list', '权限列表', '查看权限列表', 'permission', 'list', 'system', 4, 1, NOW(), NOW()),

-- 系统设置 (system 分类)
(UUID(), 'system:config', '系统配置', '修改系统配置', 'system', 'config', 'system', 0, 1, NOW(), NOW()),
(UUID(), 'system:monitor', '系统监控', '查看系统监控', 'system', 'monitor', 'system', 1, 1, NOW(), NOW()),
(UUID(), 'system:log', '查看日志', '查看系统日志', 'system', 'log', 'system', 2, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE
    display_name = VALUES(display_name),
    description = VALUES(description),
    category = VALUES(category),
    sort = VALUES(sort),
    updated_at = NOW();

-- ====================================
-- 分配权限给角色
-- ====================================

-- Admin 角色拥有所有权限
INSERT INTO role_permissions (id, role_id, permission_id, created_at)
SELECT UUID(), r.id, p.id, NOW()
FROM roles r
CROSS JOIN permissions p ON 1=1
WHERE r.name = 'admin'
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- Operator 运营人员权限
INSERT INTO role_permissions (id, role_id, permission_id, created_at)
SELECT UUID(), r.id, p.id, NOW()
FROM roles r
CROSS JOIN permissions p ON 1=1
WHERE r.name = 'operator'
AND p.name IN (
    'user:read', 'user:list',
    'order:create', 'order:read', 'order:update', 'order:list', 'order:approve', 'order:reject',
    'product:read', 'product:list',
    'inventory:read', 'inventory:list', 'inventory:update'
)
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- Merchant 商家权限
INSERT INTO role_permissions (id, role_id, permission_id, created_at)
SELECT UUID(), r.id, p.id, NOW()
FROM roles r
CROSS JOIN permissions p ON 1=1
WHERE r.name = 'merchant'
AND p.name IN (
    'order:create', 'order:read', 'order:list',
    'product:create', 'product:read', 'product:update', 'product:list',
    'inventory:create', 'inventory:read', 'inventory:update', 'inventory:list'
)
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- Customer 普通用户权限
INSERT INTO role_permissions (id, role_id, permission_id, created_at)
SELECT UUID(), r.id, p.id, NOW()
FROM roles r
CROSS JOIN permissions p ON 1=1
WHERE r.name = 'customer'
AND p.name IN (
    'order:create', 'order:read', 'order:list',
    'product:read', 'product:list'
)
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- System 系统用户权限
INSERT INTO role_permissions (id, role_id, permission_id, created_at)
SELECT UUID(), r.id, p.id, NOW()
FROM roles r
CROSS JOIN permissions p ON 1=1
WHERE r.name = 'system'
AND p.name IN (
    'product:read', 'product:update',
    'inventory:read', 'inventory:update',
    'order:read', 'order:update'
)
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- ====================================
-- 验证结果
-- ====================================
SELECT
    '权限总数' AS type,
    COUNT(*) AS count
FROM permissions
UNION ALL
SELECT
    CONCAT('角色 ', role_name, ' 权限数'),
    COUNT(*)
FROM role_permissions
GROUP BY role_name
ORDER BY type;
