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
	//user
	api.Echo.POST("/user/sign-in", api.UserHandler.HandleSignin) // xử lý yêu cầu đăng nhập người dùng
	api.Echo.POST("/user/sign-up", api.UserHandler.HandleSignup) // xử lý yêu cầu đăng ký người dùng
	// profile
	user := api.Echo.Group("/user", middlewares.JWTMiddleware())     // tạo nhóm route cho người dùng
	user.GET("/profile", api.UserHandler.HandleProfile)              // xử lý yêu cầu lấy thông tin người dùng
	user.PUT("/profile/update", api.UserHandler.HandleUpdateProfile) // xử lý yêu cầu cập nhật thông tin người dùng
}
