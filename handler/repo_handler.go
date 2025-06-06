package handler

import (
	"backend-github-trending/log"
	"backend-github-trending/model"
	"backend-github-trending/model/req"
	"backend-github-trending/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type RepoHandler struct {
	GithubRepo repository.GithubRepo
}

func NewRepoHandler(githubRepo repository.GithubRepo) *RepoHandler {
	return &RepoHandler{
		GithubRepo: githubRepo,
	}
}

// RepoTrending trả về danh sách trending repositories
func (r *RepoHandler) RepoTrending(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomClaims)

	repos, err := r.GithubRepo.SelectRepos(c.Request().Context(), claims.UserId, 25)
	if err != nil {
		log.Error("Thất bại to get trending repos:", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Thất bại to get trending repositories",
			Data:       nil,
		})
	}

	// Convert BuildBy string to Contributors array
	for i, repo := range repos {
		if repo.BuildBy != "" {
			repos[i].Contributors = strings.Split(repo.BuildBy, ",")
		}
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Successfully retrieved trending repositories",
		Data:       repos,
	})
}

// SelectBookmarks trả về danh sách bookmarks của user
func (r *RepoHandler) SelectBookmarks(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomClaims)

	repos, err := r.GithubRepo.SelectAllBookmarks(c.Request().Context(), claims.UserId)
	if err != nil {
		log.Error("Thất bại to get bookmarks:", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Thất bại to get bookmarks",
			Data:       nil,
		})
	}

	// Convert BuildBy string to Contributors array
	for i, repo := range repos {
		if repo.BuildBy != "" {
			repos[i].Contributors = strings.Split(repo.BuildBy, ",")
		}
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Successfully retrieved bookmarks",
		Data:       repos,
	})
}

// Bookmark thêm repository vào bookmark
func (r *RepoHandler) Bookmark(c echo.Context) error {
	request := req.ReqBookmark{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request format",
			Data:       nil,
		})
	}

	// Validate request
	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Validation Thất bại",
			Data:       err.Error(),
		})
	}

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomClaims)

	// Generate bookmark ID
	bookmarkID, err := uuid.NewUUID()
	if err != nil {
		log.Error("Thất bại to generate bookmark ID:", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Thất bại to generate bookmark ID",
			Data:       nil,
		})
	}

	// Add bookmark
	err = r.GithubRepo.Bookmark(c.Request().Context(), bookmarkID.String(), request.RepoName, claims.UserId)
	if err != nil {
		log.Error("Thất bại to bookmark repository:", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Thất bại to bookmark repository",
			Data:       err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Repository bookmarked successfully",
		Data:       nil,
	})
}

// DelBookmark xóa bookmark
func (r *RepoHandler) DelBookmark(c echo.Context) error {
	request := req.ReqBookmark{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request format",
			Data:       nil,
		})
	}

	// Validate request
	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Validation Thất bại",
			Data:       err.Error(),
		})
	}

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomClaims)

	// Remove bookmark
	err := r.GithubRepo.DelBookmark(c.Request().Context(), request.RepoName, claims.UserId)
	if err != nil {
		log.Error("Thất bại to remove bookmark:", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Thất bại to remove bookmark",
			Data:       err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Bookmark removed successfully",
		Data:       nil,
	})
}
