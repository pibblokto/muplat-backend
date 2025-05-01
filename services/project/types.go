package project

import "time"

type DeploymentResponse struct {
	Name      string
	Type      string
	CreatedAt time.Time
}

type ProjectResponse struct {
	Name      string
	Owner     string
	CreatedAt time.Time
}
