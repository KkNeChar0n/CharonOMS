package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// 密码迁移脚本：将明文密码转换为bcrypt加密
// 使用方法: go run scripts/migrate_passwords.go

func main() {
	// 从环境变量读取数据库配置
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "charonoms")

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	log.Println("数据库连接成功")

	// 查询所有用户
	rows, err := db.Query("SELECT id, username, password FROM useraccount")
	if err != nil {
		log.Fatalf("查询用户失败: %v", err)
	}
	defer rows.Close()

	// 统计
	var total, success, skipped, failed int

	// 准备更新语句
	updateStmt, err := db.Prepare("UPDATE useraccount SET password = ? WHERE id = ?")
	if err != nil {
		log.Fatalf("准备更新语句失败: %v", err)
	}
	defer updateStmt.Close()

	// 遍历用户
	for rows.Next() {
		var id uint
		var username, password string

		if err := rows.Scan(&id, &username, &password); err != nil {
			log.Printf("扫描用户失败: %v", err)
			failed++
			continue
		}

		total++

		// 检查是否已经是bcrypt加密（bcrypt哈希以$2a$、$2b$或$2y$开头）
		if len(password) == 60 && (password[:4] == "$2a$" || password[:4] == "$2b$" || password[:4] == "$2y$") {
			log.Printf("[跳过] 用户 %s (ID: %d) 密码已经加密", username, id)
			skipped++
			continue
		}

		// 使用bcrypt加密明文密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			log.Printf("[失败] 用户 %s (ID: %d) 密码加密失败: %v", username, id, err)
			failed++
			continue
		}

		// 更新数据库
		_, err = updateStmt.Exec(string(hashedPassword), id)
		if err != nil {
			log.Printf("[失败] 用户 %s (ID: %d) 密码更新失败: %v", username, id, err)
			failed++
			continue
		}

		log.Printf("[成功] 用户 %s (ID: %d) 密码已加密", username, id)
		success++
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("遍历用户时出错: %v", err)
	}

	// 打印统计信息
	log.Println("\n=== 密码迁移完成 ===")
	log.Printf("总用户数: %d", total)
	log.Printf("成功加密: %d", success)
	log.Printf("已跳过: %d", skipped)
	log.Printf("失败: %d", failed)

	if failed > 0 {
		log.Println("\n⚠️  有部分用户密码迁移失败，请检查日志")
		os.Exit(1)
	}

	log.Println("\n✅ 所有密码已成功迁移到bcrypt加密")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
