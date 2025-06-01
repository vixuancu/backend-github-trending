package middlewares

import (
	"backend-github-trending/model"
	req2 "backend-github-trending/model/req"
	"github.com/labstack/echo/v4"
	"net/http"
)

func IsAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// handle logic
			req := req2.ReqSignin{}
			if err := c.Bind(&req); err != nil {
				return c.JSON(http.StatusBadRequest, model.Response{
					StatusCode: http.StatusBadRequest,
					Message:    err.Error(),
					Data:       nil,
				})
			}
			// kiểm tra xem người dùng có phải là admin không
			if req.Email != "admin@gmail.com" {
				return c.JSON(http.StatusBadRequest, model.Response{
					StatusCode: http.StatusBadRequest,
					Message:    "Bạn không không có quyền gọi api này !",
					Data:       nil,
				})
			}
			return next(c)
		}
	}
}
