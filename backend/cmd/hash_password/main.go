package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("使用方法: go run main.go <密码>")
		os.Exit(1)
	}

	password := os.Args[1]

	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("生成密码哈希失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("密码: %s\n", password)
	fmt.Printf("哈希值: %s\n", hashedPassword)
}
