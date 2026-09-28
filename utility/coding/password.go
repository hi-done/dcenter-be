package coding

import "golang.org/x/crypto/bcrypt"

// HashPassword 使用 bcrypt 对密码进行单向哈希
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// VerifyPassword 校验明文密码与 bcrypt 哈希是否匹配
func VerifyPassword(password, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}
