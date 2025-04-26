package repositories

import "github.com/muplat/muplat-backend/models"

func (db *DatabaseConfig) SaveDeployment(name, projectName, deploymentType string) error {
	d := &models.Deployment{
		Name:        name,
		ProjectName: projectName,
		Type:        deploymentType,
	}
	err := db.Connection.Create(d).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DatabaseConfig) GetDeployment(name, projectName string) (*models.Deployment, error) {
	d := &models.Deployment{}
	err := db.Connection.Model(&models.Deployment{}).Where("name = ?", name).Where("project_name = ?", projectName).Take(d).Error
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (db *DatabaseConfig) DeleteDeployment(name, projectName string) error {
	d := &models.Deployment{
		Name:        name,
		ProjectName: projectName,
	}
	err := db.Connection.Model(&models.Deployment{}).Delete(d).Error
	if err != nil {
		return err
	}
	return nil
}
