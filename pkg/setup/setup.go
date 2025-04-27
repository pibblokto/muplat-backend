package setup

import (
	"github.com/muplat/muplat-backend/pkg/jwt"
	"github.com/muplat/muplat-backend/pkg/k8s"
	"github.com/muplat/muplat-backend/repositories"
)

type GlobalConfig struct {
	// Global dependencies
	Db          *repositories.Database
	ClusterConn *k8s.ClusterConnection
	Jwt         *jwt.JwtConfig
}

func InitGlobalConfig() *GlobalConfig {
	globalConf := &GlobalConfig{}
	globalConf.Db = &repositories.Database{}
	globalConf.ClusterConn = &k8s.ClusterConnection{}
	globalConf.Jwt = &jwt.JwtConfig{}

	return globalConf
}
