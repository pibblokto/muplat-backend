package models

import (
	"time"
)

type Deployment struct {
	Name           string `gorm:"primaryKey;not null;index:idx_depl,priority:2"`
	ProjectName    string `gorm:"primaryKey;not null;index:idx_proj;index:idx_depl,priority:1"`
	Type           string
	CreatedAt      time.Time
	AppConfig      AppConfig      `gorm:"foreignKey:DeploymentName,ProjectName;references:Name,ProjectName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PostgresConfig PostgresConfig `gorm:"foreignKey:DeploymentName,ProjectName;references:Name,ProjectName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
