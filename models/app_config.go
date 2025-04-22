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

func (ac *AppConfig) SaveAppConfig() (*AppConfig, error) {
	err := db.Create(ac).Error
	if err != nil {
		return &AppConfig{}, err
	}
	return ac, nil
}

func (ac *AppConfig) DeleteAppConfig() error {
	err := db.Model(&AppConfig{}).Delete(ac).Error
	if err != nil {
		return err
	}
	return nil
}
