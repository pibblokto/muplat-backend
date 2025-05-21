package repositories

import "github.com/muplat/muplat-backend/models"

func (db *Database) GetAppConfig(deploymentName, projectName string) (*models.AppConfig, error) {
	ac := &models.AppConfig{}
	err := db.Connection.
		Model(&models.AppConfig{}).
		Where("deployment_name = ?", deploymentName).
		Where("project_name = ?", projectName).
		Take(ac).Error
	if err != nil {
		return nil, err
	}
	return ac, nil
}

func (db *Database) GetAppConfigByExternalUrl(externalUrl string) (*models.AppConfig, error) {
	ac := &models.AppConfig{}
	err := db.Connection.Model(&models.AppConfig{}).Where("external_url = ?", externalUrl).Take(ac).Error
	if err != nil {
		return nil, err
	}
	return ac, nil
}

func (db *Database) SaveAppConfig(
	deploymentName,
	projectName,
	repository,
	tag,
	externalUrl,
	internalUrl,
	tier,
	envVarsSecret string,
	port uint,
) error {
	ac := &models.AppConfig{
		DeploymentName: deploymentName,
		ProjectName:    projectName,
		Repository:     repository,
		Tag:            tag,
		ExternalUrl:    externalUrl,
		InternalUrl:    internalUrl,
		Tier:           tier,
		Port:           port,
		EnvVarsSecret:  envVarsSecret,
	}
	err := db.Connection.Create(ac).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) DeleteAppConfig(deploymentName, projectName string) error {
	ac := &models.AppConfig{
		DeploymentName: deploymentName,
		ProjectName:    projectName,
	}
	err := db.Connection.Model(&models.AppConfig{}).Delete(ac).Error
	if err != nil {
		return err
	}
	return nil
}
