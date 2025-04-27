package models

type AppConfig struct {
	DeploymentName string `gorm:"primarykey; not null"`
	ProjectName    string `gorm:"primarykey; not null"`
	Repository     string
	Tag            string
	ExternalUrl    string
	InternalUrl    string
	Tier           string
	Port           uint
	EnvVarsSecret  string
}
