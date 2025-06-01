package security

import (
	"backend-github-trending/model"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Khóa bí mật dùng để ký JWT token
var secretKey = []byte("github-trending-secret-key")

// SECRET_KEY là khóa bí mật được export để sử dụng trong các package khác
var SECRET_KEY = string(secretKey)

// GetSecretKey trả về khóa bí mật dạng []byte để sử dụng trong JWT
func GetSecretKey() []byte {
	return secretKey
}

// GenToken tạo JWT token cho người dùng
func GenToken(user model.User) (string, error) {
	// Tạo claims với thông tin người dùng
	claims := model.JwtCustomClaims{
		UserId: user.UserId,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Tạo token với thuật toán HS256 và claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Ký token với khóa bí mật
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken kiểm tra tính hợp lệ của token và trả về claims
func ValidateToken(tokenString string) (*model.JwtCustomClaims, error) {
	// Parse token với key function
	token, err := jwt.ParseWithClaims(tokenString, &model.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra thuật toán ký
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Kiểm tra token có hợp lệ không
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Chuyển đổi claims sang JwtCustomClaims
	claims, ok := token.Claims.(*model.JwtCustomClaims)
	if !ok {
		return nil, errors.New("failed to parse claims")
	}

	return claims, nil
}

// ExtractTokenFromHeader trích xuất token từ header Authorization
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("authorization header format must be 'Bearer {token}'")
	}
	return authHeader[7:], nil
}
