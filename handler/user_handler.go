package handler

import (
	"backend-github-trending/log"
	"backend-github-trending/model"
	req2 "backend-github-trending/model/req"
	"backend-github-trending/repository"
	"backend-github-trending/security"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserHandler struct {
	UserRepo repository.UserRepo
}

func (u *UserHandler) HandleSignup(c echo.Context) error {
	req := req2.ReqSignup{}
	if err := c.Bind(&req); err != nil {
		return err // Middleware sẽ xử lý lỗi này
	}

	if err := c.Validate(&req); err != nil {
		return err // Middleware sẽ xử lý lỗi này
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
	// tạo token
	token, err := security.GenToken(user)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, model.Response{
			http.StatusInternalServerError,
			"Lỗi khi tạo token",
			nil,
		})
	}
	user.Token = token
	return c.JSON(http.StatusOK, model.Response{
		http.StatusOK,
		"Đăng ký thành công",
		user,
	})
}

func (u *UserHandler) HandleSignin(c echo.Context) error {
	req := req2.ReqSignin{}
	if err := c.Bind(&req); err != nil { // lấy dữ liệu từ request body
		log.Error(err.Error())
		// Kiểm tra lỗi EOF cụ thể - thường xảy ra khi request body rỗng
		if err.Error() == "EOF" {
			return c.JSON(http.StatusBadRequest, model.Response{
				StatusCode: http.StatusBadRequest,
				Message:    "Request body trống. Vui lòng cung cấp email và mật khẩu.",
				Data:       nil,
			})
		}
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Định dạng request không hợp lệ: " + err.Error(),
			Data:       nil,
		})
	}

	// Kiểm tra thêm các trường trống sau khi binding
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Email và mật khẩu là bắt buộc",
			Data:       nil,
		})
	}

	// Sử dụng custom validator thay vì validator mặc định
	if err := c.Validate(&req); err != nil {
		// ValidationMiddleware sẽ xử lý
		return err
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
	// tạo token
	token, err := security.GenToken(user)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, model.Response{
			http.StatusInternalServerError,
			"Lỗi khi tạo token",
			nil,
		})
	}
	user.Token = token // gán token cho user
	return c.JSON(http.StatusOK, model.Response{
		http.StatusOK,
		"Đăng nhập thành công",
		user,
	})
}

// HandleProfile xử lý yêu cầu lấy thông tin người dùng
func (u *UserHandler) HandleProfile(c echo.Context) error {
	// Lấy thông tin user từ token JWT
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)

	// Lấy userId từ claims
	userId := claims.UserId

	// Lấy thông tin người dùng từ database
	userInfo, err := u.UserRepo.FindByID(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Lỗi khi lấy thông tin người dùng",
			Data:       nil,
		})
	}

	// Không trả về password cho client
	userInfo.Password = ""

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Lấy thông tin người dùng thành công",
		Data:       userInfo,
	})
}
func (u *UserHandler) HandleUpdateProfile(c echo.Context) error {
	// Lấy thông tin user từ token JWT
	jwtToken := c.Get("user").(*jwt.Token)
	claims := jwtToken.Claims.(*model.JwtCustomClaims)

	// Lấy userId từ claims
	userId := claims.UserId

	// Bind request data
	req := req2.RequestUpdateUser{}
	if err := c.Bind(&req); err != nil {
		return err // Middleware sẽ xử lý lỗi này
	}

	// Validate request data
	if err := c.Validate(&req); err != nil {
		return err // Middleware sẽ xử lý lỗi này
	}

	// Tạo đối tượng User mới với thông tin cập nhật
	updatedUser := &model.User{
		UserId:   userId,
		Fullname: req.Fullname,
		Email:    req.Email,
	}

	// Cập nhật thông tin người dùng
	userInfo, err := u.UserRepo.UpdateProfile(c.Request().Context(), updatedUser)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	// Không trả về password cho client
	userInfo.Password = ""

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Cập nhật thông tin người dùng thành công",
		Data:       userInfo,
	})
}
