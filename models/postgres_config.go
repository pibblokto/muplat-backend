package models

type PostgresConfig struct {
	DeploymentName    string `gorm:"primarykey; not null"`
	ProjectName       string `gorm:"primarykey; not null"`
	DiskSize          string
	InternalEndpoint  string
	Database          string
	CredentialsSecret string
}
