package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = bcrypt.DefaultCost

// HashPassword 对明文密码做 bcrypt 哈希，用于存储
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash 校验明文密码与存储的哈希是否一致
// 若 hash 以 "$2" 开头则按 bcrypt 校验，否则按明文比对（兼容旧数据迁移）
func CheckPasswordHash(password, hash string) bool {
	if len(hash) >= 2 && hash[:2] == "$2" {
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		return err == nil
	}
	return password == hash
}
