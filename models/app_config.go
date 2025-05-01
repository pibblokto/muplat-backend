package models

type AppConfig struct {
	DeploymentName string `gorm:"primarykey; not null; index:aconf"`
	ProjectName    string `gorm:"primarykey; not null; index:aconf"`
	Repository     string
	Tag            string
	ExternalUrl    string
	InternalUrl    string
	Tier           string
	Port           uint
	EnvVarsSecret  string
}
