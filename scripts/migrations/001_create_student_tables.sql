-- Migration Script: Create Student Tables
-- Date: 2026-01-27
-- Description: Create student and student_coach tables for student management module

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE charonoms;

-- Create student table
CREATE TABLE IF NOT EXISTS `student` (
  `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '学生ID',
  `name` VARCHAR(100) NOT NULL COMMENT '学生姓名',
  `sex_id` INT NOT NULL COMMENT '性别ID',
  `grade_id` INT NOT NULL COMMENT '年级ID',
  `phone` VARCHAR(20) NOT NULL COMMENT '联系电话',
  `status` TINYINT DEFAULT 0 COMMENT '状态（0=启用，1=禁用）',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  INDEX `idx_sex_id` (`sex_id`),
  INDEX `idx_grade_id` (`grade_id`),
  INDEX `idx_status` (`status`),

  CONSTRAINT `fk_student_sex` FOREIGN KEY (`sex_id`) REFERENCES `sex` (`id`),
  CONSTRAINT `fk_student_grade` FOREIGN KEY (`grade_id`) REFERENCES `grade` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学生表';

-- Create student_coach relationship table
CREATE TABLE IF NOT EXISTS `student_coach` (
  `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '关联ID',
  `student_id` INT NOT NULL COMMENT '学生ID',
  `coach_id` INT NOT NULL COMMENT '教练ID',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

  UNIQUE KEY `uk_student_coach` (`student_id`, `coach_id`),
  INDEX `idx_student_id` (`student_id`),
  INDEX `idx_coach_id` (`coach_id`),

  CONSTRAINT `fk_sc_student` FOREIGN KEY (`student_id`) REFERENCES `student` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_sc_coach` FOREIGN KEY (`coach_id`) REFERENCES `coach` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学生教练关联表';

-- Verification queries
SELECT 'Student tables created successfully!' AS status;
SELECT 'student table structure:' AS info;
DESCRIBE student;
SELECT 'student_coach table structure:' AS info;
DESCRIBE student_coach;
