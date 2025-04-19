package setup

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type MuplatCfg struct {
	ConnectionMode   string `env:"CONNECTION_MODE" envDefault:"internal"`
	KubeconfigPath   string `env:"KUBECONFIG"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresPort     string `env:"POSTGRES_PORT" envDefault:"5432"`
	Database         string `env:"DATABASE"`
	InitUser         string `env:"INIT_USER" envDefault:"admin"`
	InitUserPassword string `env:"INIT_USER_PASSWORD"`
	JwtLifespanHours string `env:"JWT_LIFESPAN_HOURS" envDefault:"1"`
	JwtSecret        string `env:"JWT_SECRET" envDefault:"1"`
}

func LoadConfig() MuplatCfg {
	var conf MuplatCfg

	err := env.Parse(&conf)
	if err != nil {
		log.Fatal(err.Error())
		return MuplatCfg{}
	}
	if conf.ConnectionMode != "internal" && conf.ConnectionMode != "external" {
		log.Fatal("CONNECTION_MODE should be either internal or external")
		return MuplatCfg{}
	}
	if conf.ConnectionMode == "external" && conf.KubeconfigPath == "" {
		log.Fatalf("KUBECONFIG is required if CONNECTION_MODE is \"%s\"", conf.ConnectionMode)
		return MuplatCfg{}
	}

	return conf
}
