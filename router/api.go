package router

import (
	"backend-github-trending/handler"
	"backend-github-trending/middlewares"
	"github.com/labstack/echo/v4"
)

type API struct {
	Echo        *echo.Echo           // làm việc với echo framework
	UserHandler handler.UserHandler  // xử lý các yêu cầu liên quan đến người dùng
	RepoHandler *handler.RepoHandler // Thêm RepoHandler
}

func (api *API) SetupRouter() {
	//user
	api.Echo.POST("/user/sign-in", api.UserHandler.HandleSignin) // xử lý yêu cầu đăng nhập người dùng
	api.Echo.POST("/user/sign-up", api.UserHandler.HandleSignup) // xử lý yêu cầu đăng ký người dùng
	// profile
	user := api.Echo.Group("/user", middlewares.JWTMiddleware())     // tạo nhóm route cho người dùng
	user.GET("/profile", api.UserHandler.HandleProfile)              // xử lý yêu cầu lấy thông tin người dùng
	user.PUT("/profile/update", api.UserHandler.HandleUpdateProfile) // xử lý yêu cầu cập nhật thông tin người dùng

	// github repo
	github := api.Echo.Group("/github", middlewares.JWTMiddleware())
	github.GET("/trending", api.RepoHandler.RepoTrending)
	github.POST("/trending", api.RepoHandler.RepoTrending) // lưu trữ repo trending

	// bookmark
	bookmark := api.Echo.Group("/bookmark", middlewares.JWTMiddleware())
	bookmark.GET("/list", api.RepoHandler.SelectBookmarks)
	bookmark.POST("/add", api.RepoHandler.Bookmark)
	bookmark.DELETE("/delete", api.RepoHandler.DelBookmark)
}
