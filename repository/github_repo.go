package repository

import (
	"backend-github-trending/model"
	"context"
)

// Interface định nghĩa các method mà implementation phải có
type GithubRepo interface {
	// Repository operations
	SaveRepo(context context.Context, repo model.GithubRepo) (model.GithubRepo, error)
	SelectRepos(context context.Context, userId string, limit int) ([]model.GithubRepo, error)
	SelectRepoByName(context context.Context, name string) (model.GithubRepo, error)
	UpdateRepo(context context.Context, repo model.GithubRepo) (model.GithubRepo, error)

	// Bookmark operations
	SelectAllBookmarks(context context.Context, userId string) ([]model.GithubRepo, error)
	Bookmark(context context.Context, bid, nameRepo, userId string) error
	DelBookmark(context context.Context, nameRepo, userId string) error
}

// Tại sao dùng Interface?
// - Dễ test (có thể tạo mock)
// - Dễ thay đổi implementation (từ PostgreSQL sang MongoDB chẳng hạn)
// - Code clean hơn, tách biệt concerns
