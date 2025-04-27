package deployment

import "fmt"

type DeploymentType string
type AppTier string
type PostgresTier string

var (
	// deployment types
	TypeApp      DeploymentType = "app"
	TypePostgres DeploymentType = "postgres"
	// app tier
	AppTierDev AppTier = "dev"
	AppTierMid AppTier = "mid"
	AppTierPro AppTier = "pro"
)

func (dt DeploymentType) IsValid() bool {
	switch dt {
	case TypeApp, TypePostgres:
		return true
	default:
		return false
	}
}

func (at AppTier) IsValid() bool {
	switch at {
	case AppTierDev, AppTierMid, AppTierPro:
		return true
	default:
		return false
	}
}

func (at *AppTier) UnmarshalJSON(b []byte) error {
	s := string(b[1 : len(b)-1])
	temp := AppTier(s)
	if !temp.IsValid() {
		return fmt.Errorf("invalid app tier type: %s", s)
	}
	*at = temp
	return nil
}

func (dt *DeploymentType) UnmarshalJSON(b []byte) error {
	s := string(b[1 : len(b)-1])
	temp := DeploymentType(s)
	if !temp.IsValid() {
		return fmt.Errorf("invalid deployment type: %s", s)
	}
	*dt = temp
	return nil
}

type AppConfig struct {
	Repository string            `json:"repository" binding:"required"`
	Tag        string            `json:"tag" binding:"required"`
	External   *bool             `json:"external" binding:"required"`
	DomainName string            `json:"domainName"`
	Tier       AppTier           `json:"tier" binding:"required"`
	Port       uint              `json:"port" binding:"required"`
	EnvVars    map[string]string `json:"envVars"`
}

type PostgresConfig struct {
	DiskSize uint    `json:"diskSize" binding:"required"`
	Database *string `json:"database"`
}
