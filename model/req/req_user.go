package req

type RequestUpdateUser struct {
	Fullname string `json:"fullname" validate:"required,min=3,max=50"` // Tên đầy đủ của người dùng, bắt buộc, từ 3 đến 50 ký tự
	Email    string `json:"email" validate:"required,email"`           // Email của người dùng, bắt buộc và phải là định dạng email hợp lệ
}
