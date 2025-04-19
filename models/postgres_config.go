package models

type PostgresConfig struct {
	DeploymentName    string `gorm:"primarykey; not null"`
	ProjectName       string `gorm:"primarykey; not null"`
	DiskSize          string
	Endpoint          string
	Database          string
	Tier              string
	CredentialsSecret string
}

func (pc *PostgresConfig) SavePostgresConfig() (*PostgresConfig, error) {
	err := db.Create(pc).Error
	if err != nil {
		return &PostgresConfig{}, err
	}
	return pc, nil
}
