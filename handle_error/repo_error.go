package handle_error

import "errors"

var (
	// Repository errors
	RepoNotFound   = errors.New("repository not found")
	RepoConflict   = errors.New("repository already exists")
	RepoInsertFail = errors.New("failed to insert repository")
	RepoNotUpdated = errors.New("repository not updated")

	// Bookmark errors
	BookmarkNotFound = errors.New("bookmark not found")
	BookmarkConflict = errors.New("bookmark already exists")
	BookmarkFail     = errors.New("failed to create bookmark")
	DelBookmarkFail  = errors.New("failed to delete bookmark")
)
