-- 初始化基础数据：性别、年级、学科

-- 创建性别表
CREATE TABLE IF NOT EXISTS `sex` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(10) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='性别表';

-- 插入性别数据
INSERT INTO `sex` (`id`, `name`) VALUES
(1, '男'),
(2, '女')
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`);

-- 创建年级表
CREATE TABLE IF NOT EXISTS `grade` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0-启用 1-禁用',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='年级表';

-- 插入年级数据
INSERT INTO `grade` (`id`, `name`, `status`) VALUES
(1, '一年级', 0),
(2, '二年级', 0),
(3, '三年级', 0),
(4, '四年级', 0),
(5, '五年级', 0),
(6, '六年级', 0),
(7, '初一', 0),
(8, '初二', 0),
(9, '初三', 0),
(10, '高一', 0),
(11, '高二', 0),
(12, '高三', 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `status`=VALUES(`status`);

-- 创建学科表
CREATE TABLE IF NOT EXISTS `subject` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `subject` varchar(50) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0-启用 1-禁用',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学科表';

-- 插入学科数据
INSERT INTO `subject` (`id`, `subject`, `status`) VALUES
(1, '语文', 0),
(2, '数学', 0),
(3, '英语', 0),
(4, '物理', 0),
(5, '化学', 0),
(6, '生物', 0),
(7, '历史', 0),
(8, '地理', 0),
(9, '政治', 0),
(10, '音乐', 0),
(11, '美术', 0),
(12, '体育', 0)
ON DUPLICATE KEY UPDATE `subject`=VALUES(`subject`), `status`=VALUES(`status`);
