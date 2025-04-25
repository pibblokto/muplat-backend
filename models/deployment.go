package models

import (
	"time"
)

type Deployment struct {
	Name           string `gorm:"primarykey; not null"`
	ProjectName    string `gorm:"primarykey; not null"`
	CreatedAt      time.Time
	Type           string
	AppConfig      AppConfig      `gorm:"foreignKey:DeploymentName,ProjectName;references:Name,ProjectName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PostgresConfig PostgresConfig `gorm:"foreignKey:DeploymentName,ProjectName;references:Name,ProjectName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (d *Deployment) GetDeployment(name, project string) error {
	err := db.Model(&Deployment{}).Where("name = ?", name).Where("project_name = ?", project).Take(d).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *Deployment) SaveDeployment() (*Deployment, error) {
	err := db.Create(d).Error
	if err != nil {
		return &Deployment{}, err
	}
	return d, nil
}

func (d *Deployment) DeleteDeployment() error {
	err := db.Model(&Deployment{}).Delete(d).Error
	if err != nil {
		return err
	}
	return nil
}
