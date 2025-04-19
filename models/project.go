package models

import "time"

type Project struct {
	Name          string `gorm:"primarykey"`
	CreatedAt     time.Time
	Owner         string
	Namespace     string
	NetworkPolicy string
	Deployments   []Deployment `gorm:"foreignKey:ProjectName;references:Name;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (p *Project) SaveProject() (*Project, error) {
	err := db.Create(p).Error
	if err != nil {
		return &Project{}, err
	}
	return p, nil
}
