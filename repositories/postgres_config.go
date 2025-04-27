package repositories

import "github.com/muplat/muplat-backend/models"

func (db *Database) SavePostgresConfig(
	deploymentName,
	projectName,
	diskSize,
	internalEndpoint,
	database,
	credentialsSecret string,
) error {
	pc := &models.PostgresConfig{
		DeploymentName:    deploymentName,
		ProjectName:       projectName,
		DiskSize:          diskSize,
		InternalEndpoint:  internalEndpoint,
		Database:          database,
		CredentialsSecret: credentialsSecret,
	}

	err := db.Connection.Create(pc).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) DeletePostgresConfig(
	deploymentName,
	projectName string,
) error {
	pc := &models.PostgresConfig{
		DeploymentName: deploymentName,
		ProjectName:    projectName,
	}
	err := db.Connection.Model(&models.PostgresConfig{}).Delete(pc).Error
	if err != nil {
		return err
	}
	return nil
}
