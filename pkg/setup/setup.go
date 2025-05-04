package setup

import (
	"github.com/muplat/muplat-backend/pkg/jwt"
	"github.com/muplat/muplat-backend/pkg/k8s"
	"github.com/muplat/muplat-backend/repositories"
)

type GlobalConfig struct {
	// Global dependencies
	PlatformDomain string `env:"PLATFORM_DOMAIN"`
	Db             *repositories.Database
	ClusterConn    *k8s.ClusterConnection
	Jwt            *jwt.JwtConfig
}

func InitGlobalConfig() *GlobalConfig {
	globalConf := &GlobalConfig{}
	globalConf.Db = repositories.NewDatabase()
	globalConf.ClusterConn = k8s.NewClusterConnection()
	globalConf.Jwt = jwt.NewJwtConfig()
	return globalConf
}
