package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 生成管理员密码哈希
	password := "admin123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("================================")
	fmt.Println("管理员密码哈希生成工具")
	fmt.Println("================================")
	fmt.Printf("原始密码: %s\n", password)
	fmt.Printf("哈希值: %s\n", string(hash))
	fmt.Println("\n请将此哈希值插入到数据库：")
	fmt.Printf("INSERT INTO admins (username, password, role) VALUES ('admin', '%s', 'admin');\n", string(hash))
	fmt.Println("================================")
}
