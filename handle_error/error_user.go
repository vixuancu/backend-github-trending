package handle_error

import "errors"

var (
	UserNotFound      = errors.New("Không có User này")
	UserConflig       = errors.New("User đã tồn tại trong hệ thống")
	SignupFail        = errors.New("Đăng kí thất bại")
	UpdateProfileFail = errors.New("Cập nhật thông tin người dùng thất bại")
)
