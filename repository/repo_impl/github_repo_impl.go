package repo_impl

import (
	"backend-github-trending/db"
	"backend-github-trending/handle_error"
	"backend-github-trending/log"
	"backend-github-trending/model"
	"backend-github-trending/repository"
	"context"
	"database/sql"
	"github.com/lib/pq"
	"time"
)

type GithubRepoImpl struct {
	sql *db.Sql
}

func NewGithubRepo(sql *db.Sql) repository.GithubRepo {
	return &GithubRepoImpl{
		sql: sql,
	}
}

// SaveRepo lưu repository mới vào database
func (g *GithubRepoImpl) SaveRepo(context context.Context, repo model.GithubRepo) (model.GithubRepo, error) {
	statement := `INSERT INTO repos(
					name, description, url, color, lang, fork, stars, 
 			        stars_today, build_by, created_at, updated_at) 
          		  VALUES(
					:name, :description, :url, :color, :lang, :fork, :stars, 
					:stars_today, :build_by, :created_at, :updated_at
				  )`

	repo.CreatedAt = time.Now()
	repo.UpdatedAt = time.Now()

	_, err := g.sql.Db.NamedExecContext(context, statement, repo)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return repo, handle_error.RepoConflict
		}
		log.Error(err.Error())
		return repo, handle_error.RepoInsertFail
	}

	return repo, nil
}

// SelectRepos lấy danh sách repositories với bookmark status
func (g *GithubRepoImpl) SelectRepos(context context.Context, userId string, limit int) ([]model.GithubRepo, error) {
	var repos []model.GithubRepo
	err := g.sql.Db.SelectContext(context, &repos,
		`SELECT 
			repos.name, repos.description, repos.url, repos.color, repos.lang, 
			repos.fork, repos.stars, repos.stars_today, repos.build_by, 
			repos.created_at, repos.updated_at,
			COALESCE(repos.name = bookmarks.repo_name, FALSE) as bookmarked
		FROM repos
		LEFT JOIN bookmarks 
		ON repos.name = bookmarks.repo_name AND bookmarks.user_id = $1  
		WHERE repos.name IS NOT NULL 
		ORDER BY repos.updated_at DESC LIMIT $2`, userId, limit)

	if err != nil {
		log.Error(err.Error())
		return repos, err
	}

	return repos, nil
}

// SelectRepoByName tìm repository theo tên
func (g *GithubRepoImpl) SelectRepoByName(context context.Context, name string) (model.GithubRepo, error) {
	var repo model.GithubRepo
	err := g.sql.Db.GetContext(context, &repo,
		`SELECT * FROM repos WHERE name = $1`, name)

	if err != nil {
		if err == sql.ErrNoRows {
			return repo, handle_error.RepoNotFound
		}
		log.Error(err.Error())
		return repo, err
	}
	return repo, nil
}

// UpdateRepo cập nhật thông tin repository
func (g *GithubRepoImpl) UpdateRepo(context context.Context, repo model.GithubRepo) (model.GithubRepo, error) {
	sqlStatement := `
		UPDATE repos
		SET 
			stars = :stars,
			fork = :fork,
			stars_today = :stars_today,
			build_by = :build_by,
			updated_at = :updated_at
		WHERE name = :name`

	repo.UpdatedAt = time.Now()

	result, err := g.sql.Db.NamedExecContext(context, sqlStatement, repo)
	if err != nil {
		log.Error(err.Error())
		return repo, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
		return repo, handle_error.RepoNotUpdated
	}
	if count == 0 {
		return repo, handle_error.RepoNotUpdated
	}

	return repo, nil
}

// SelectAllBookmarks lấy tất cả bookmarks của user
func (g *GithubRepoImpl) SelectAllBookmarks(context context.Context, userId string) ([]model.GithubRepo, error) {
	var repos []model.GithubRepo
	err := g.sql.Db.SelectContext(context, &repos,
		`SELECT 
			repos.name, repos.description, repos.url, 
			repos.color, repos.lang, repos.fork, repos.stars, 
			repos.stars_today, repos.build_by, repos.created_at, repos.updated_at,
			true as bookmarked
		FROM bookmarks 
		INNER JOIN repos ON repos.name = bookmarks.repo_name
		WHERE bookmarks.user_id = $1
		ORDER BY bookmarks.created_at DESC`, userId)

	if err != nil {
		if err == sql.ErrNoRows {
			return repos, handle_error.BookmarkNotFound
		}
		log.Error(err.Error())
		return repos, err
	}
	return repos, nil
}

// Bookmark thêm bookmark cho user
func (g *GithubRepoImpl) Bookmark(context context.Context, bid, nameRepo, userId string) error {
	statement := `INSERT INTO bookmarks(
					bid, user_id, repo_name, created_at, updated_at) 
          		  VALUES($1, $2, $3, $4, $5)`

	now := time.Now()
	_, err := g.sql.Db.ExecContext(context, statement, bid, userId, nameRepo, now, now)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return handle_error.BookmarkConflict
		}
		log.Error(err.Error())
		return handle_error.BookmarkFail
	}

	return nil
}

// DelBookmark xóa bookmark của user
func (g *GithubRepoImpl) DelBookmark(context context.Context, nameRepo, userId string) error {
	result, err := g.sql.Db.ExecContext(context,
		"DELETE FROM bookmarks WHERE repo_name = $1 AND user_id = $2",
		nameRepo, userId)

	if err != nil {
		log.Error(err.Error())
		return handle_error.DelBookmarkFail
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
		return handle_error.DelBookmarkFail
	}

	if rowsAffected == 0 {
		return handle_error.BookmarkNotFound
	}

	return nil
}
