package models

type PostgresConfig struct {
	DeploymentName    string `gorm:"primarykey; not null;"`
	ProjectName       string `gorm:"primarykey; not null;index:idx_proj"`
	DiskSize          string
	InternalEndpoint  string
	Database          string
	CredentialsSecret string
}
