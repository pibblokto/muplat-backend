package models

type PostgresConfig struct {
	DeploymentName    string `gorm:"primarykey; not null;index:pconf"`
	ProjectName       string `gorm:"primarykey; not null;index:pconf"`
	DiskSize          string
	InternalEndpoint  string
	Database          string
	CredentialsSecret string
}
