package user

import "time"

type ProjectResponse struct {
	Name      string
	Owner     string
	CreatedAt time.Time
}

type UserResponse struct {
	Name      string
	IsAdmin   bool
	CreatedAt time.Time
}
