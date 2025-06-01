package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value"`
}

type ValidationResponse struct {
	Errors []ValidationError `json:"errors"`
}

// CustomValidator wrapper cho validator
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator tạo instance mới của CustomValidator
func NewValidator() *CustomValidator {
	v := validator.New()

	// Sử dụng tên JSON tag thay vì tên struct field
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Đăng ký các custom validators
	registerCustomValidators(v)

	return &CustomValidator{
		validator: v,
	}
}

// Validate implements echo.Validator để sử dụng với Echo framework
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, cv.ValidateStruct(i))
	}
	return nil
}

// ValidateStruct kiểm tra validation và trả về lỗi thân thiện
func (cv *CustomValidator) ValidateStruct(s interface{}) *ValidationResponse {
	err := cv.validator.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors []ValidationError

	for _, err := range err.(validator.ValidationErrors) {
		// Lấy tên field từ json tag nếu có, nếu không dùng tên field
		fieldName := getFieldName(s, err)

		validationErrors = append(validationErrors, ValidationError{
			Field:   fieldName,
			Message: getErrorMessage(err),
			Value:   fmt.Sprintf("%v", err.Value()),
		})
	}

	return &ValidationResponse{
		Errors: validationErrors,
	}
}

// getFieldName lấy tên field từ json tag
func getFieldName(s interface{}, err validator.FieldError) string {
	field := err.Field()

	// Cố gắng lấy tên từ json tag
	if t := reflect.TypeOf(s); t.Kind() == reflect.Struct {
		if f, exists := t.FieldByName(field); exists {
			if jsonTag := f.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
				field = strings.Split(jsonTag, ",")[0]
			}
		}
	}

	return strings.ToLower(field)
}

// getErrorMessage trả về thông báo lỗi thân thiện
func getErrorMessage(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())
	vietnameseName := getVietnameseFieldName(field)

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s là bắt buộc", vietnameseName)
	case "email":
		return fmt.Sprintf("%s phải có định dạng email hợp lệ", vietnameseName)
	case "min":
		if fe.Kind().String() == "string" {
			return fmt.Sprintf("%s phải có ít nhất %s ký tự", vietnameseName, fe.Param())
		}
		return fmt.Sprintf("%s phải có giá trị tối thiểu là %s", vietnameseName, fe.Param())
	case "max":
		if fe.Kind().String() == "string" {
			return fmt.Sprintf("%s không được vượt quá %s ký tự", vietnameseName, fe.Param())
		}
		return fmt.Sprintf("%s không được vượt quá %s", vietnameseName, fe.Param())
	case "len":
		return fmt.Sprintf("%s phải có độ dài chính xác %s ký tự", vietnameseName, fe.Param())
	case "gte":
		return fmt.Sprintf("%s phải lớn hơn hoặc bằng %s", vietnameseName, fe.Param())
	case "lte":
		return fmt.Sprintf("%s phải nhỏ hơn hoặc bằng %s", vietnameseName, fe.Param())
	case "alphanum":
		return fmt.Sprintf("%s chỉ được chứa chữ cái và số", vietnameseName)
	case "alpha":
		return fmt.Sprintf("%s chỉ được chứa chữ cái", vietnameseName)
	case "numeric":
		return fmt.Sprintf("%s chỉ được chứa số", vietnameseName)
	case "oneof":
		return fmt.Sprintf("%s phải là một trong các giá trị: %s", vietnameseName, fe.Param())
	case "url":
		return fmt.Sprintf("%s phải là URL hợp lệ", vietnameseName)
	case "uuid":
		return fmt.Sprintf("%s phải là UUID hợp lệ", vietnameseName)
	case "eqfield":
		return fmt.Sprintf("%s phải giống với %s", vietnameseName, getVietnameseFieldName(fe.Param()))
	case "nefield":
		return fmt.Sprintf("%s không được giống với %s", vietnameseName, getVietnameseFieldName(fe.Param()))
	case "gt":
		return fmt.Sprintf("%s phải lớn hơn %s", vietnameseName, fe.Param())
	case "lt":
		return fmt.Sprintf("%s phải nhỏ hơn %s", vietnameseName, fe.Param())
	case "containsany":
		return fmt.Sprintf("%s phải chứa ít nhất một trong các ký tự: %s", vietnameseName, fe.Param())
	case "containsrune":
		return fmt.Sprintf("%s phải chứa ký tự: %s", vietnameseName, fe.Param())
	case "excludesall":
		return fmt.Sprintf("%s không được chứa bất kỳ ký tự nào trong: %s", vietnameseName, fe.Param())
	case "excludesrune":
		return fmt.Sprintf("%s không được chứa ký tự: %s", vietnameseName, fe.Param())
	case "startswith":
		return fmt.Sprintf("%s phải bắt đầu bằng: %s", vietnameseName, fe.Param())
	case "endswith":
		return fmt.Sprintf("%s phải kết thúc bằng: %s", vietnameseName, fe.Param())
	case "password":
		return fmt.Sprintf("%s phải có ít nhất một chữ hoa, một chữ thường, một chữ số và một ký tự đặc biệt", vietnameseName)
	default:
		return fmt.Sprintf("%s không hợp lệ", vietnameseName)
	}
}

// registerCustomValidators đăng ký các custom validators
func registerCustomValidators(v *validator.Validate) {
	// Validator cho mật khẩu mạnh (ít nhất 1 chữ hoa, 1 chữ thường, 1 số và 1 ký tự đặc biệt)
	v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		hasUpper := false
		hasLower := false
		hasNumber := false
		hasSpecial := false

		for _, c := range value {
			if c >= 'A' && c <= 'Z' {
				hasUpper = true
			} else if c >= 'a' && c <= 'z' {
				hasLower = true
			} else if c >= '0' && c <= '9' {
				hasNumber = true
			} else {
				hasSpecial = true
			}
		}

		return hasUpper && hasLower && hasNumber && hasSpecial
	})

	// Thêm các custom validators khác ở đây
}

// ValidationMiddleware tạo middleware để tự động validate request body
func ValidationMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Gọi next handler và check errors
			err := next(c)

			// Nếu là validation error, trả về thông báo lỗi thân thiện
			if he, ok := err.(*echo.HTTPError); ok && he.Code == http.StatusBadRequest {
				if ve, ok := he.Message.(*ValidationResponse); ok {
					return c.JSON(http.StatusBadRequest, map[string]interface{}{
						"status_code": http.StatusBadRequest,
						"message":     "Dữ liệu không hợp lệ",
						"errors":      ve.Errors,
					})
				}
			}

			return err
		}
	}
}

// getVietnameseFieldName chuyển đổi tên field sang tiếng Việt
func getVietnameseFieldName(field string) string {
	fieldMap := map[string]string{
		"fullname":         "Họ và tên",
		"email":            "Email",
		"password":         "Mật khẩu",
		"phone":            "Số điện thoại",
		"address":          "Địa chỉ",
		"age":              "Tuổi",
		"username":         "Tên đăng nhập",
		"firstname":        "Tên",
		"lastname":         "Họ",
		"dob":              "Ngày sinh",
		"confirm_password": "Xác nhận mật khẩu",
		"role":             "Vai trò",
		"title":            "Tiêu đề",
		"content":          "Nội dung",
		"avatar":           "Ảnh đại diện",
		"description":      "Mô tả",
		"status":           "Trạng thái",
		"created_at":       "Ngày tạo",
		"updated_at":       "Ngày cập nhật",
		"deleted_at":       "Ngày xóa",
		"id":               "ID",
		"userid":           "ID người dùng",
		"user_id":          "ID người dùng",
		"url":              "URL",
		"code":             "Mã",
		"price":            "Giá",
		"quantity":         "Số lượng",
		"name":             "Tên",
	}

	if vietnameseName, exists := fieldMap[field]; exists {
		return vietnameseName
	}
	return strings.ReplaceAll(field, "_", " ") // Format snake_case thành "snake case"
}
