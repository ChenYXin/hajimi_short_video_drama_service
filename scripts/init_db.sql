-- 数据库初始化脚本
-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS hajimi CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE hajimi;

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    avatar VARCHAR(255) DEFAULT '',
    phone VARCHAR(20) DEFAULT '',
    status ENUM('active', 'inactive', 'banned') DEFAULT 'active',
    last_login_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建管理员表
CREATE TABLE IF NOT EXISTS admins (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role ENUM('admin', 'super_admin', 'editor') DEFAULT 'admin',
    permissions JSON,
    status ENUM('active', 'inactive') DEFAULT 'active',
    last_login_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_role (role),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建短剧表
CREATE TABLE IF NOT EXISTS dramas (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    cover_image VARCHAR(255) DEFAULT '',
    category VARCHAR(50) DEFAULT '',
    tags JSON,
    director VARCHAR(100) DEFAULT '',
    actors JSON,
    release_date DATE,
    status ENUM('draft', 'published', 'archived') DEFAULT 'draft',
    view_count BIGINT UNSIGNED DEFAULT 0,
    like_count BIGINT UNSIGNED DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0.00,
    duration INT UNSIGNED DEFAULT 0, -- 总时长（秒）
    episode_count INT UNSIGNED DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_title (title),
    INDEX idx_category (category),
    INDEX idx_status (status),
    INDEX idx_release_date (release_date),
    INDEX idx_view_count (view_count),
    INDEX idx_created_at (created_at),
    FULLTEXT idx_search (title, description)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建剧集表
CREATE TABLE IF NOT EXISTS episodes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    drama_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    episode_num INT UNSIGNED NOT NULL,
    video_url VARCHAR(500) DEFAULT '',
    thumbnail VARCHAR(255) DEFAULT '',
    duration INT UNSIGNED DEFAULT 0, -- 时长（秒）
    status ENUM('draft', 'published', 'archived') DEFAULT 'draft',
    view_count BIGINT UNSIGNED DEFAULT 0,
    like_count BIGINT UNSIGNED DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (drama_id) REFERENCES dramas(id) ON DELETE CASCADE,
    INDEX idx_drama_id (drama_id),
    INDEX idx_episode_num (episode_num),
    INDEX idx_status (status),
    INDEX idx_view_count (view_count),
    INDEX idx_created_at (created_at),
    UNIQUE KEY uk_drama_episode (drama_id, episode_num)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建用户观看历史表
CREATE TABLE IF NOT EXISTS user_watch_history (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    drama_id BIGINT UNSIGNED NOT NULL,
    episode_id BIGINT UNSIGNED NOT NULL,
    watch_progress INT UNSIGNED DEFAULT 0, -- 观看进度（秒）
    watch_duration INT UNSIGNED DEFAULT 0, -- 观看时长（秒）
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (drama_id) REFERENCES dramas(id) ON DELETE CASCADE,
    FOREIGN KEY (episode_id) REFERENCES episodes(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_drama_id (drama_id),
    INDEX idx_episode_id (episode_id),
    INDEX idx_created_at (created_at),
    UNIQUE KEY uk_user_episode (user_id, episode_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建用户收藏表
CREATE TABLE IF NOT EXISTS user_favorites (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    drama_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (drama_id) REFERENCES dramas(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_drama_id (drama_id),
    INDEX idx_created_at (created_at),
    UNIQUE KEY uk_user_drama (user_id, drama_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建评论表
CREATE TABLE IF NOT EXISTS comments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    drama_id BIGINT UNSIGNED NOT NULL,
    episode_id BIGINT UNSIGNED NULL,
    content TEXT NOT NULL,
    rating TINYINT UNSIGNED DEFAULT 0, -- 1-5星评分
    like_count BIGINT UNSIGNED DEFAULT 0,
    status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (drama_id) REFERENCES dramas(id) ON DELETE CASCADE,
    FOREIGN KEY (episode_id) REFERENCES episodes(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_drama_id (drama_id),
    INDEX idx_episode_id (episode_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    config_key VARCHAR(100) NOT NULL UNIQUE,
    config_value TEXT,
    description VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_config_key (config_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入默认系统配置
INSERT INTO system_configs (config_key, config_value, description) VALUES
('site_name', 'Gin MySQL API', '网站名称'),
('site_description', '基于Gin和MySQL的短剧API系统', '网站描述'),
('upload_max_size', '100', '文件上传最大大小(MB)'),
('video_allowed_types', '["mp4", "avi", "mov", "mkv", "webm"]', '允许的视频文件类型'),
('image_allowed_types', '["jpg", "jpeg", "png", "gif", "webp"]', '允许的图片文件类型'),
('cache_ttl', '3600', '缓存过期时间(秒)'),
('pagination_limit', '20', '分页默认限制'),
('max_pagination_limit', '100', '分页最大限制')
ON DUPLICATE KEY UPDATE 
    config_value = VALUES(config_value),
    updated_at = CURRENT_TIMESTAMP;

-- 创建默认超级管理员账户
-- 密码: admin123 (BCrypt 哈希)
INSERT INTO admins (username, email, password, role, status) VALUES
('admin', 'admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'super_admin', 'active')
ON DUPLICATE KEY UPDATE 
    password = VALUES(password),
    role = VALUES(role),
    status = VALUES(status),
    updated_at = CURRENT_TIMESTAMP;

-- 创建触发器：更新短剧的剧集数量
DELIMITER $$

CREATE TRIGGER IF NOT EXISTS update_drama_episode_count_insert
AFTER INSERT ON episodes
FOR EACH ROW
BEGIN
    UPDATE dramas 
    SET episode_count = (
        SELECT COUNT(*) 
        FROM episodes 
        WHERE drama_id = NEW.drama_id AND deleted_at IS NULL
    )
    WHERE id = NEW.drama_id;
END$$

CREATE TRIGGER IF NOT EXISTS update_drama_episode_count_delete
AFTER UPDATE ON episodes
FOR EACH ROW
BEGIN
    IF NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL THEN
        UPDATE dramas 
        SET episode_count = (
            SELECT COUNT(*) 
            FROM episodes 
            WHERE drama_id = NEW.drama_id AND deleted_at IS NULL
        )
        WHERE id = NEW.drama_id;
    END IF;
END$$

DELIMITER ;

-- 创建视图：热门短剧
CREATE OR REPLACE VIEW popular_dramas AS
SELECT 
    d.*,
    COALESCE(AVG(c.rating), 0) as avg_rating,
    COUNT(DISTINCT c.id) as comment_count,
    COUNT(DISTINCT f.id) as favorite_count
FROM dramas d
LEFT JOIN comments c ON d.id = c.drama_id AND c.status = 'approved' AND c.deleted_at IS NULL
LEFT JOIN user_favorites f ON d.id = f.drama_id
WHERE d.status = 'published' AND d.deleted_at IS NULL
GROUP BY d.id
ORDER BY d.view_count DESC, d.like_count DESC;

-- 创建视图：用户统计
CREATE OR REPLACE VIEW user_stats AS
SELECT 
    u.id,
    u.username,
    u.email,
    u.status,
    u.created_at,
    COUNT(DISTINCT f.drama_id) as favorite_count,
    COUNT(DISTINCT h.drama_id) as watched_drama_count,
    COUNT(DISTINCT c.id) as comment_count,
    MAX(h.updated_at) as last_watch_time
FROM users u
LEFT JOIN user_favorites f ON u.id = f.user_id
LEFT JOIN user_watch_history h ON u.id = h.user_id
LEFT JOIN comments c ON u.id = c.user_id AND c.deleted_at IS NULL
WHERE u.deleted_at IS NULL
GROUP BY u.id;

-- 创建存储过程：清理过期数据
DELIMITER $$

CREATE PROCEDURE IF NOT EXISTS CleanupExpiredData()
BEGIN
    DECLARE done INT DEFAULT FALSE;
    DECLARE cleanup_date DATE DEFAULT DATE_SUB(CURDATE(), INTERVAL 90 DAY);
    
    -- 清理90天前的观看历史（保留最近观看记录）
    DELETE h1 FROM user_watch_history h1
    INNER JOIN (
        SELECT user_id, episode_id, MIN(id) as keep_id
        FROM user_watch_history
        WHERE created_at < cleanup_date
        GROUP BY user_id, episode_id
    ) h2 ON h1.user_id = h2.user_id AND h1.episode_id = h2.episode_id
    WHERE h1.id != h2.keep_id AND h1.created_at < cleanup_date;
    
    -- 清理已删除数据的软删除记录（超过30天）
    DELETE FROM users WHERE deleted_at IS NOT NULL AND deleted_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
    DELETE FROM dramas WHERE deleted_at IS NOT NULL AND deleted_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
    DELETE FROM episodes WHERE deleted_at IS NOT NULL AND deleted_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
    DELETE FROM comments WHERE deleted_at IS NOT NULL AND deleted_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
    
    SELECT 'Cleanup completed' as result;
END$$

DELIMITER ;

-- 创建事件调度器（每天凌晨2点执行清理）
-- SET GLOBAL event_scheduler = ON;
-- CREATE EVENT IF NOT EXISTS daily_cleanup
-- ON SCHEDULE EVERY 1 DAY STARTS '2024-01-01 02:00:00'
-- DO CALL CleanupExpiredData();

COMMIT;