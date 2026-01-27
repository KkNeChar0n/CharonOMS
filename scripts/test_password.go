package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// 密码加密测试脚本
// 用于验证bcrypt加密和验证功能
// 使用方法: go run scripts/test_password.go

func main() {
	fmt.Println("========================================")
	fmt.Println("密码加密功能测试")
	fmt.Println("========================================")
	fmt.Println()

	// 测试密码
	testPasswords := []string{
		"password",
		"admin123",
		"Test@123",
		"复杂密码123!@#",
	}

	for i, password := range testPasswords {
		fmt.Printf("测试 #%d: 密码 = \"%s\"\n", i+1, password)
		fmt.Println("----------------------------------------")

		// 加密密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			log.Printf("❌ 加密失败: %v\n\n", err)
			continue
		}

		fmt.Printf("✅ 加密成功\n")
		fmt.Printf("   明文密码: %s\n", password)
		fmt.Printf("   加密后: %s\n", string(hashedPassword))
		fmt.Printf("   长度: %d 字符\n", len(hashedPassword))

		// 验证正确密码
		err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
		if err != nil {
			fmt.Printf("❌ 验证失败 (应该成功): %v\n", err)
		} else {
			fmt.Printf("✅ 正确密码验证成功\n")
		}

		// 验证错误密码
		wrongPassword := password + "_wrong"
		err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(wrongPassword))
		if err != nil {
			fmt.Printf("✅ 错误密码验证失败 (符合预期)\n")
		} else {
			fmt.Printf("❌ 错误密码验证成功 (不应该发生)\n")
		}

		fmt.Println()
	}

	fmt.Println("========================================")
	fmt.Println("测试完成")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("说明:")
	fmt.Println("- bcrypt哈希长度为60字符")
	fmt.Println("- 格式: $2a$10$<salt><hash>")
	fmt.Println("- 每次加密相同密码会产生不同的哈希（随机salt）")
	fmt.Println("- 验证时会自动提取salt进行比对")
}
