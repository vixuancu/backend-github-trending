package middlewares

import (
	"backend-github-trending/model"
	"backend-github-trending/security"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware trả về một Echo middleware function sử dụng JWT để xác thực
func JWTMiddleware() echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &model.JwtCustomClaims{}
		},
		SigningKey: security.GetSecretKey(),
	}

	return echojwt.WithConfig(config)
}
