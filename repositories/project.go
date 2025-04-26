package repositories

import "github.com/muplat/muplat-backend/models"

func (db *DatabaseConfig) SaveProject(name, owner, namespace, networkPolicy string) error {
	p := &models.Project{
		Name:          name,
		Owner:         owner,
		Namespace:     namespace,
		NetworkPolicy: networkPolicy,
	}
	err := db.Connection.Create(p).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DatabaseConfig) GetPorjectByName(name string) (*models.Project, error) {
	p := &models.Project{}
	err := db.Connection.Model(&models.Project{}).Where("name = ?", name).Take(p).Error
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (db *DatabaseConfig) DeleteProject(name string) error {
	p := &models.Project{
		Name: name,
	}
	err := db.Connection.Model(&models.Project{}).Delete(p).Error
	if err != nil {
		return err
	}
	return nil
}
