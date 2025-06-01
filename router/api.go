package router

import (
	"backend-github-trending/handler"
	"backend-github-trending/middlewares"
	"github.com/labstack/echo/v4"
)

type API struct {
	Echo        *echo.Echo          // làm việc với echo framework
	UserHandler handler.UserHandler // xử lý các yêu cầu liên quan đến người dùng
}

func (api *API) SetupRoter() {

	api.Echo.POST("/user/sign-in", api.UserHandler.HandleSignin) // xử lý yêu cầu đăng nhập người dùng
	api.Echo.POST("/user/sign-up", api.UserHandler.HandleSignup) // xử lý yêu cầu đăng ký người dùng
	api.Echo.POST("/user/profile", api.UserHandler.HandleProfile, middlewares.JWTMiddleware())
}
