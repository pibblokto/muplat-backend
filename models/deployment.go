package models

import (
	"time"
)

type Deployment struct {
	Name           string `gorm:"primarykey; not null"`
	ProjectName    string `gorm:"primarykey; not null"`
	Type           string
	CreatedAt      time.Time
	AppConfig      AppConfig      `gorm:"foreignKey:DeploymentName,ProjectName;references:Name,ProjectName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PostgresConfig PostgresConfig `gorm:"foreignKey:DeploymentName,ProjectName;references:Name,ProjectName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
