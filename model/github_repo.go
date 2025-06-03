package model

import "time"

type GithubRepo struct {
	Name         string    `json:"name" db:"name"`               // Tên repo (vd: "microsoft/typescript")
	Url          string    `json:"url" db:"url"`                 // Link đến repo
	Description  string    `json:"description" db:"description"` // Mô tả repo
	Color        string    `json:"color" db:"color"`             // Màu của ngôn ngữ (vd: "#007ACC")
	Lang         string    `json:"lang" db:"lang"`               // Ngôn ngữ chính (vd: "TypeScript")
	Fork         string    `json:"fork" db:"fork"`               // Số fork
	Stars        string    `json:"stars" db:"stars"`             // Số stars
	StarsToday   string    `json:"stars_today" db:"stars_today"` // Stars hôm nay
	BuildBy      string    `json:"-" db:"build_by"`
	Bookmarked   bool      `json:"bookmarked" db:"bookmarked"` // User đã bookmark chưa
	Contributors []string  `json:"contributors,omitempty"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Bookmark model
type Bookmark struct {
	BID       string    `json:"bid" db:"bid"`
	UserID    string    `json:"user_id" db:"user_id"`
	RepoName  string    `json:"repo_name" db:"repo_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Giải thích các tag:
// - `json:"name"` : Tên field khi convert sang JSON
// - `db:"name"`   : Tên column trong database
// - `json:"-"`    : Không hiển thị field này trong JSON response
// - `omitempty`   : Nếu field rỗng thì không hiển thị trong JSON
