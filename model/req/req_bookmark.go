package req

type ReqBookmark struct {
	RepoName string `json:"repo_name" validate:"required" `
}
