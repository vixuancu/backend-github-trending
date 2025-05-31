package handler

import (
	"backend-github-trending/model"
	req2 "backend-github-trending/model/req"
	"backend-github-trending/repository"
	"backend-github-trending/security"
	"github.com/go-playground/validator/v10"
	uuid "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

type UserHandler struct {
	UserRepo repository.UserRepo
}

func (u *UserHandler) HandleSignup(c echo.Context) error {
	req := req2.ReqSignup{}
	if err := c.Bind(&req); err != nil { // lấy dữ liệu từ request body
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			http.StatusBadRequest,
			err.Error(),
			nil,
		})
	}
	// validate dữ liệu bằng thư viện validator
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	// mã hóa mật khẩu
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			http.StatusInternalServerError,
			"lỗi khi mã hóa mật khẩu",
			err.Error(),
		})
	}
	role := model.MEMBER.String()
	// uuid
	userId, err := uuid.NewUUID() // tạo userId bằng uuid
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			http.StatusInternalServerError,
			"lỗi khi tạo userId",
			err.Error(),
		})
	}
	user := model.User{
		UserId:   userId.String(),
		Fullname: req.Fullname,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     role,
		Token:    "",
	}
	user, err = u.UserRepo.SaveUser(c.Request().Context(), &user) // lưu user vào cơ sở dữ liệu
	if err != nil {
		log.Error(err.Error())

		return c.JSON(http.StatusConflict, model.Response{
			http.StatusConflict,
			err.Error(),
			nil,
		})

	}
	return c.JSON(http.StatusOK, model.Response{
		http.StatusOK,
		"Đăng ký thành công",
		nil,
	})
}

func (u *UserHandler) HandleSignin(c echo.Context) error {
	req := req2.ReqSignin{}
	if err := c.Bind(&req); err != nil { // lấy dữ liệu từ request body
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			http.StatusBadRequest,
			err.Error(),
			nil,
		})
	}
	// validate dữ liệu bằng thư viện validator
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	user, err := u.UserRepo.CheckLogin(c.Request().Context(), req) // kiểm tra đăng nhập
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.Response{
			http.StatusUnauthorized,
			err.Error(),
			nil,
		})
	}
	// check pass
	if !security.CheckPasswordHash(req.Password, user.Password) {
		log.Error("Mật khẩu không đúng")
		return c.JSON(http.StatusUnauthorized, model.Response{
			http.StatusUnauthorized,
			"Mật khẩu không đúng",
			nil,
		})
	}
	return c.JSON(http.StatusOK, model.Response{
		http.StatusOK,
		"Đăng nhập thành công",
		user,
	})
}
