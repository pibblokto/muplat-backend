package repositories

import "github.com/muplat/muplat-backend/models"

func (db *DatabaseConfig) SaveAppConfig(
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

func (db *DatabaseConfig) DeleteAppConfig(deploymentName, projectName string) error {
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
