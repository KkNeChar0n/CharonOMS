package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 连接数据库
	dsn := "root:qweasd123Q!@tcp(localhost:3306)/charonoms?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}
	defer db.Close()

	// 获取admin用户密码
	var username, password string
	var userID int
	err = db.QueryRow("SELECT id, username, password FROM useraccount WHERE username = 'admin'").Scan(&userID, &username, &password)
	if err != nil {
		log.Fatalf("查询用户失败: %v", err)
	}

	fmt.Printf("用户: %s (ID: %d)\n", username, userID)
	fmt.Printf("密码Hash: %s\n\n", password)

	// 测试不同的原始密码
	testPasswords := []string{
		"password",
		"admin123",
		"admin",
		"123456",
	}

	fmt.Println("测试原始密码:")
	fmt.Println("----------------------------------------")

	for _, testPwd := range testPasswords {
		err = bcrypt.CompareHashAndPassword([]byte(password), []byte(testPwd))
		if err == nil {
			fmt.Printf("✅ '%s' - 匹配成功!\n", testPwd)
		} else {
			fmt.Printf("❌ '%s' - 不匹配\n", testPwd)
		}
	}
}
