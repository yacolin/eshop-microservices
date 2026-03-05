-- ====================================
-- RBAC 数据迁移脚本（向后兼容）
-- ====================================
-- 用途: 从旧格式迁移到新的 RBAC 格式
-- 使用方式: mysql -u root -p user_db < scripts/rbac-migration.sql
-- 注意: 此脚本应谨慎使用，建议先备份数据库
-- ====================================

USE user_db;

-- ====================================
-- 检查并创建必要的表（如果不存在）
-- ====================================

-- 创建 roles 表（如果不存在）
CREATE TABLE IF NOT EXISTS roles (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    status TINYINT DEFAULT 1,
    sort INT DEFAULT 0,
    is_system TINYINT(1) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_name (name),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建 permissions 表（如果不存在）
CREATE TABLE IF NOT EXISTS permissions (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    category VARCHAR(50) NOT NULL,
    sort INT DEFAULT 0,
    status TINYINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_name (name),
    INDEX idx_resource_action (resource, action),
    INDEX idx_category (category),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建 role_permissions 表（如果不存在）
CREATE TABLE IF NOT EXISTS role_permissions (
    id VARCHAR(36) PRIMARY KEY,
    role_id VARCHAR(36) NOT NULL,
    permission_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE KEY uk_role_permission (role_id, permission_id, deleted_at),
    INDEX idx_role_id (role_id),
    INDEX idx_permission_id (permission_id),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建 user_roles 表（如果不存在）
CREATE TABLE IF NOT EXISTS user_roles (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    role_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE KEY uk_user_role (user_id, role_id, deleted_at),
    INDEX idx_user_id (user_id),
    INDEX idx_role_id (role_id),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ====================================
-- 迁移旧数据（如果存在）
-- ====================================

-- 检查是否存在旧的 user_roles 表结构（带有 role_name 字段）
SET @old_table_exists = (
    SELECT COUNT(*)
    FROM information_schema.columns
    WHERE table_schema = 'user_db'
    AND table_name = 'user_roles'
    AND column_name = 'role_name'
);

-- 如果存在旧的 user_roles 表结构，进行数据迁移
SET @sql = IF(@old_table_exists > 0,
    CONCAT('
        -- 备份旧数据到临时表
        CREATE TEMPORARY TABLE temp_old_user_roles AS
        SELECT user_id, role_name FROM user_roles WHERE deleted_at IS NULL;
        
        -- 删除旧表中的 role_name 字段数据
        DELETE FROM user_roles WHERE role_name IS NOT NULL;
        
        -- 迁移数据：将 role_name 转换为 role_id
        INSERT INTO user_roles (id, user_id, role_id, created_at)
        SELECT 
            UUID(),
            t.user_id,
            r.id,
            NOW()
        FROM temp_old_user_roles t
        JOIN roles r ON t.role_name = r.name
        ON DUPLICATE KEY UPDATE role_id = r.id, updated_at = NOW();
        
        -- 删除临时表
        DROP TEMPORARY TABLE IF EXISTS temp_old_user_roles;
    '),
    'SELECT ''No old table structure found, skipping migration'' AS message;'
);

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ====================================
-- 数据完整性检查
-- ====================================

-- 检查孤立的角色权限记录（角色或权限已被删除）
SELECT
    '孤立的角色权限记录' AS check_type,
    COUNT(*) AS count
FROM role_permissions rp
LEFT JOIN roles r ON rp.role_id = r.id
LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE r.id IS NULL OR p.id IS NULL;

-- 检查孤立的用户角色记录（用户或角色已被删除）
SELECT
    '孤立的用户角色记录' AS check_type,
    COUNT(*) AS count
FROM user_roles ur
LEFT JOIN users u ON ur.user_id = u.id
LEFT JOIN roles r ON ur.role_id = r.id
WHERE u.id IS NULL OR r.id IS NULL;

-- 清理孤立记录（可选，取消注释以执行）
-- DELETE FROM role_permissions WHERE role_id NOT IN (SELECT id FROM roles);
-- DELETE FROM role_permissions WHERE permission_id NOT IN (SELECT id FROM permissions);
-- DELETE FROM user_roles WHERE user_id NOT IN (SELECT id FROM users);
-- DELETE FROM user_roles WHERE role_id NOT IN (SELECT id FROM roles);

-- ====================================
-- 验证迁移结果
-- ====================================
SELECT
    '迁移完成验证' AS section,
    '角色数' AS item,
    COUNT(*) AS value
FROM roles
WHERE deleted_at IS NULL
UNION ALL
SELECT
    '迁移完成验证',
    '权限数',
    COUNT(*)
FROM permissions
WHERE deleted_at IS NULL
UNION ALL
SELECT
    '迁移完成验证',
    '角色权限映射数',
    COUNT(*)
FROM role_permissions
WHERE deleted_at IS NULL
UNION ALL
SELECT
    '迁移完成验证',
    '用户角色映射数',
    COUNT(*)
FROM user_roles
WHERE deleted_at IS NULL;