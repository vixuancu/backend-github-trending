package handle_error

import "errors"

var (
	UserNotFound = errors.New("Không tìm thấy User")
	UserConflig  = errors.New("User đã tồn tại trong hệ thống")
	SignupFail   = errors.New("Đăng kí thất bại")
)
