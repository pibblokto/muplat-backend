setup  package

import (
	"github.com/muplat/muplat-backend/pkg/jwt"
	"github.com/muplat/muplat-backend/pkg/k8s"
	"github.com/muplat/muplat-backend/repositories"
	"github.com/muplat/muplat-backend/services"
)

type GlobalConfig struct {
	// Global dependencies
	Db  *repositories.DatabaseConfig
	K8s *k8s.ClusterConfig
	Jwt *jwt.JwtConfig
}

func InitGlobalConfig() *GlobalConfig {
	platform := &GlobalConfig{}
	platform.Db = &repositories.DatabaseConfig{}
	platform.K8s = &k8s.ClusterConfig{}
	platform.Jwt = &jwt.JwtConfig{}

	// Initialize global dependencies
	platform.Db.InitDatabase()
	platform.K8s.InitClusterConnection()
	platform.Jwt.InitJwt()

	return platform
}
