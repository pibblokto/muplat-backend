package models

import "time"

type Project struct {
	Name          string `gorm:"primarykey"`
	CreatedAt     time.Time
	Owner         string
	Namespace     string
	NetworkPolicy string
	Deployments   []Deployment `gorm:"foreignKey:ProjectName;references:Name;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
