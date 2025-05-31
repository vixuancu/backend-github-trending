package model

import "time"

type User struct {
	UserId    string    `json:"user_id,omitempty" db:"user_id"`
	Fullname  string    `json:"fullname,omitempty" db:"full_name"`
	Email     string    `json:"email,omitempty" db:"email"`
	Password  string    `json:"password,omitempty" db:"password"`
	Role      string    `json:"role,omitempty" db:"role"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
	Token     string    `json:"token,omitempty" db:"token"`
}
