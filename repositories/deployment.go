package repositories

import "github.com/muplat/muplat-backend/models"

func (db *Database) SaveDeployment(name, projectName, deploymentType string) error {
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

func (db *Database) GetDeployment(name, projectName string) (*models.Deployment, error) {
	d := &models.Deployment{}
	err := db.Connection.Model(&models.Deployment{}).Where("name = ?", name).Where("project_name = ?", projectName).Take(d).Error
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (db *Database) DeleteDeployment(name, projectName string) error {
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

func (db *Database) GetDeployments() ([]*models.Deployment, error) {
	d := []*models.Deployment{}
	err := db.Connection.Model(&models.Deployment{}).Order("created_at DESC").Find(&d).Error

	if err != nil {
		return nil, err
	}

	return d, err
}

func (db *Database) GetDeploymentsByProject(projectName string) ([]*models.Deployment, error) {
	d := []*models.Deployment{}
	err := db.Connection.Model(&models.Deployment{}).Where("project_name = ?", projectName).Order("created_at DESC").Find(&d).Error
	if err != nil {
		return nil, err
	}
	return d, nil
}
