package repository

import (
	"backend-github-trending/model"
	"backend-github-trending/model/req"
	"context"
)

type UserRepo interface {
	SaveUser(context context.Context, user *model.User) (model.User, error)
	CheckLogin(context context.Context, loginReq req.ReqSignin) (model.User, error)
	FindByID(userId string) (*model.User, error) // Tìm kiếm người dùng theo ID
}
