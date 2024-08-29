
CREATE TABLE `tbl_file` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `file_hash` CHAR(64) NOT NULL DEFAULT '' COMMENT '文件hash', -- SHA-256
  `file_name` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `file_size` BIGINT DEFAULT 0 COMMENT '文件大小',
  `file_path` VARCHAR(512) NOT NULL DEFAULT '' COMMENT '文件存储位置',
  `create_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `update_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期',
  `status` ENUM('active', 'disabled', 'deleted') NOT NULL DEFAULT 'active' COMMENT '状态',
  `reference_count` INT DEFAULT 0 COMMENT '文件引用计数',
  `ext1` INT DEFAULT 0 COMMENT '备用字段1',
  `ext2` TEXT COMMENT '备用字段2',
  UNIQUE KEY `idx_file_hash` (`file_hash`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `tbl_user` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `username` VARCHAR(64) NOT NULL UNIQUE COMMENT '用户名',
  `password` VARCHAR(60) NOT NULL DEFAULT '' COMMENT '用户encoded密码',
  `email` VARCHAR(64) DEFAULT '' COMMENT '邮箱',
  `phone` VARCHAR(20) DEFAULT '' COMMENT '手机号',
  `email_validated` TINYINT(1) DEFAULT 0 COMMENT '邮箱是否已验证',
  `phone_validated` TINYINT(1) DEFAULT 0 COMMENT '手机号是否已验证',
  `signup_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '注册日期',
  `last_active` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活跃时间戳',
  `profile` JSON COMMENT '用户属性', -- 使用 JSON 数据类型
  `status` ENUM('active', 'disabled', 'locked', 'deleted') NOT NULL DEFAULT 'active' COMMENT '账户状态',
  UNIQUE KEY `idx_username` (`username`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `tbl_user_file` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `user_id` INT NOT NULL COMMENT '用户ID',
  `file_id` INT NOT NULL COMMENT '文件ID',
  `file_name` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `upload_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  `status` ENUM('active', 'disabled', 'deleted') NOT NULL DEFAULT 'active' COMMENT '文件状态',
  FOREIGN KEY (`user_id`) REFERENCES `tbl_user`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`file_id`) REFERENCES `tbl_file`(`id`) ON DELETE CASCADE,
  UNIQUE KEY `idx_user_file` (`user_id`, `file_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;