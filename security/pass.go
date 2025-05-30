package security

import (
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

// mã hóa mật khẩu bằng bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	return string(bytes), nil
}

// so sánh mật khẩu với hash đã mã hóa
// Returns true if they match, false otherwise
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Debugf("password comparison failed: %v", err)
		return false
	}
	return true
}
