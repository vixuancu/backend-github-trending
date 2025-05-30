package model

import "time"

type User struct {
	UserId    string    `json:"user_id,omitempty"`    // ID người dùng omitempty // không trả về nếu là 0
	Fullname  string    `json:"fullname,omitempty"`   // Tên đầy đủ người dùng
	Email     string    `json:"email,omitempty"`      // Email người dùng
	Password  string    `json:"password,omitempty"`   // Mật khẩu người dùng
	Role      string    `json:"role,omitempty"`       // Vai trò người dùng (ví dụ: admin, user)
	CreatedAt time.Time `json:"created_at,omitempty"` // Ngày tạo tài khoản
	UpdatedAt time.Time `json:"updated_at,omitempty"` // Ngày cập nhật tài khoản
	Token     string    `json:"token,omitempty"`      // Token xác thực người dùng
}
