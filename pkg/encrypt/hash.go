package encrypt

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

// Md5 计算给定字节切片的MD5散列值并返回十六进制字符串表示。
// 该函数用于需要快速、简单地生成数据的唯一标识时，但不适用于安全敏感场景。
func Md5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

// GenPasswordHash 使用bcrypt算法为密码生成哈希值。
// 该函数用于安全地存储密码，确保即使数据库泄露，密码也不会轻易被破解。
// 参数:
//
//	password []byte: 待哈希的原始密码。
//
// 返回值:
//
//	[]byte: 生成的密码哈希值。
//	error: 如果生成哈希值过程中发生错误，返回错误信息。
func GenPasswordHash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// ValidatePasswordHash 验证给定密码与存储的哈希值是否匹配。
// 该函数用于登录验证等场景，确保用户输入的密码与之前存储的哈希值一致。
// 参数:
//
//	password string: 用户输入的原始密码。
//	hashed string: 存储的密码哈希值。
//
// 返回值:
//
//	bool: 如果密码与哈希值匹配返回true，否则返回false。
func ValidatePasswordHash(password string, hashed string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return false
	}
	return true
}
