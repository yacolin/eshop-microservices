-- ====================================
-- 用户角色初始化脚本
-- ====================================
-- 用途: 为现有用户分配默认角色
-- 使用方式: mysql -u root -p user_db < scripts/user-roles-init.sql
-- ====================================

USE user_db;

-- ====================================
-- 为现有用户分配默认角色
-- ====================================

-- 为所有现有用户分配 'customer' 角色（如果还没有角色）
INSERT INTO user_roles (id, user_id, role_id, created_at)
SELECT 
    UUID(),
    u.id,
    (SELECT id FROM roles WHERE name = 'customer' LIMIT 1),
    NOW()
FROM users u
WHERE u.deleted_at IS NULL
AND NOT EXISTS (
    SELECT 1 FROM user_roles ur WHERE ur.user_id = u.id
)
ON DUPLICATE KEY UPDATE created_at = NOW();

-- ====================================
-- 创建示例管理员用户（如果不存在）
-- ====================================

-- 检查是否已存在管理员用户
SET @admin_exists = (SELECT COUNT(*) FROM users WHERE id = 'admin-user-id');

-- 如果不存在，创建管理员用户
INSERT INTO users (id, status, created_at, updated_at)
SELECT 'admin-user-id', 1, NOW(), NOW()
WHERE @admin_exists = 0;

-- 创建管理员用户信息
INSERT INTO user_infos (id, user_id, nickname, avatar, created_at, updated_at)
SELECT UUID(), 'admin-user-id', '系统管理员', 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin', NOW(), NOW()
WHERE @admin_exists = 0;

-- 为管理员用户分配 admin 角色
INSERT INTO user_roles (id, user_id, role_id, created_at)
SELECT 
    UUID(),
    'admin-user-id',
    (SELECT id FROM roles WHERE name = 'admin' LIMIT 1),
    NOW()
WHERE @admin_exists = 0;

-- ====================================
-- 创建示例商家用户（如果不存在）
-- ====================================

-- 检查是否已存在商家用户
SET @merchant_exists = (SELECT COUNT(*) FROM users WHERE id = 'merchant-user-id');

-- 如果不存在，创建商家用户
INSERT INTO users (id, status, created_at, updated_at)
SELECT 'merchant-user-id', 1, NOW(), NOW()
WHERE @merchant_exists = 0;

-- 创建商家用户信息
INSERT INTO user_infos (id, user_id, nickname, avatar, created_at, updated_at)
SELECT UUID(), 'merchant-user-id', '示例商家', 'https://api.dicebear.com/7.x/avataaars/svg?seed=merchant', NOW(), NOW()
WHERE @merchant_exists = 0;

-- 为商家用户分配 merchant 角色
INSERT INTO user_roles (id, user_id, role_id, created_at)
SELECT 
    UUID(),
    'merchant-user-id',
    (SELECT id FROM roles WHERE name = 'merchant' LIMIT 1),
    NOW()
WHERE @merchant_exists = 0;

-- ====================================
-- 创建示例运营人员用户（如果不存在）
-- ====================================

-- 检查是否已存在运营用户
SET @operator_exists = (SELECT COUNT(*) FROM users WHERE id = 'operator-user-id');

-- 如果不存在，创建运营用户
INSERT INTO users (id, status, created_at, updated_at)
SELECT 'operator-user-id', 1, NOW(), NOW()
WHERE @operator_exists = 0;

-- 创建运营用户信息
INSERT INTO user_infos (id, user_id, nickname, avatar, created_at, updated_at)
SELECT UUID(), 'operator-user-id', '运营人员', 'https://api.dicebear.com/7.x/avataaars/svg?seed=operator', NOW(), NOW()
WHERE @operator_exists = 0;

-- 为运营用户分配 operator 角色
INSERT INTO user_roles (id, user_id, role_id, created_at)
SELECT 
    UUID(),
    'operator-user-id',
    (SELECT id FROM roles WHERE name = 'operator' LIMIT 1),
    NOW()
WHERE @operator_exists = 0;

-- ====================================
-- 验证结果
-- ====================================
SELECT
    '用户总数' AS type,
    COUNT(*) AS count
FROM users
WHERE deleted_at IS NULL
UNION ALL
SELECT
    '有角色的用户数',
    COUNT(DISTINCT user_id)
FROM user_roles
WHERE deleted_at IS NULL
UNION ALL
SELECT
    CONCAT('角色 ', r.name, ' 用户数'),
    COUNT(DISTINCT ur.user_id)
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
WHERE ur.deleted_at IS NULL
GROUP BY r.name
ORDER BY type;