package models

import "time"

type Project struct {
	Name          string `gorm:"primarykey;uniqueIndex"`
	Owner         string
	Namespace     string
	NetworkPolicy string
	CreatedAt     time.Time
	Deployments   []Deployment `gorm:"foreignKey:ProjectName;references:Name;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
