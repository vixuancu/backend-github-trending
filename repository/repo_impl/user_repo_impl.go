package repo_impl

import (
	"backend-github-trending/db"
	"backend-github-trending/handle_error"
	"backend-github-trending/model"
	"backend-github-trending/model/req"
	"backend-github-trending/repository"
	"context"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	"time"
)

type UserRepoImpl struct {
	sql *db.Sql
}

func NewUserRepoImpl(sql *db.Sql) repository.UserRepo {
	return &UserRepoImpl{
		sql: sql,
	}
}

// tìm kiếm user trong cơ sở dữ liệu
func (u *UserRepoImpl) FindUserByEmail(context context.Context, email string) (model.User, error) {
	statement := `SELECT user_id, full_name, email, password, role, created_at, updated_at FROM users WHERE email = $1`
	var user model.User
	err := u.sql.Db.GetContext(context, &user, statement, email)
	if err != nil {
		log.Error(err.Error())
		return model.User{}, err
	}
	return user, nil
}

// lưu user bằng postgreQl
func (u *UserRepoImpl) SaveUser(context context.Context, user *model.User) (model.User, error) {
	// Check if the user already exists
	existingUser, err := u.FindUserByEmail(context, user.Email)
	if err == nil && existingUser.UserId != "" { // User exists
		return *user, handle_error.UserConflig
	}

	if err != nil && err.Error() != "sql: no rows in result set" { // Handle unexpected errors
		log.Error(err.Error())
		return *user, handle_error.SignupFail
	}

	// Prepare to insert the new user
	statement := `INSERT INTO users (user_id, full_name, email, password, role, created_at, updated_at) 
			VALUES (:user_id, :full_name, :email, :password, :role, :created_at, :updated_at)`
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = u.sql.Db.NamedExecContext(context, statement, user)
	if err != nil {
		log.Error(err.Error())
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return *user, handle_error.UserConflig
		}
		return *user, handle_error.SignupFail
	}

	return *user, nil
}

func (u *UserRepoImpl) CheckLogin(context context.Context, loginReq req.ReqSignin) (model.User, error) {
	user, err := u.FindUserByEmail(context, loginReq.Email)
	if err != nil {
		// So sánh chuỗi lỗi thay vì so sánh đối tượng lỗi
		if err.Error() == "sql: no rows in result set" {
			return user, handle_error.UserNotFound
		}
		log.Error(err.Error())
		return user, err
	}

	// Trả về user nếu tìm thấy (việc xác minh mật khẩu sẽ được thực hiện ở handler)
	return user, nil
}

// FindByID tìm kiếm người dùng theo ID
func (u *UserRepoImpl) FindByID(userId string) (*model.User, error) {
	statement := `SELECT user_id, full_name, email, password, role, created_at, updated_at FROM users WHERE user_id = $1`
	var user model.User
	err := u.sql.Db.Get(&user, statement, userId)
	if err != nil {
		log.Error(err.Error())
		if err.Error() == "sql: no rows in result set" {
			return nil, handle_error.UserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// UpdateProfile cập nhật thông tin người dùng
func (u *UserRepoImpl) UpdateProfile(context context.Context, user *model.User) (model.User, error) {
	statement := `UPDATE users SET full_name = :full_name, email = :email, updated_at = :updated_at WHERE user_id = :user_id`
	user.UpdatedAt = time.Now()

	_, err := u.sql.Db.NamedExecContext(context, statement, user)
	if err != nil {
		log.Error(err.Error())
		return *user, handle_error.UpdateProfileFail
	}

	// Trả về người dùng đã cập nhật
	return *user, nil
}
