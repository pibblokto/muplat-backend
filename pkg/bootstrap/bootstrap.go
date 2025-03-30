package bootstrap

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type AppCfg struct {
	ConnectionMode string `env:"CONNECTION_MODE" envDefault:"internal"`
	KubeconfigPath string `env:"KUBECONFIG"`
}

func LoadConfig() AppCfg {
	var conf AppCfg

	err := env.Parse(&conf)
	log.Print(conf.ConnectionMode)
	log.Print(conf.KubeconfigPath)
	if err != nil {
		log.Fatal(err.Error())
		return AppCfg{}
	}
	if conf.ConnectionMode != "internal" && conf.ConnectionMode != "external" {
		log.Fatal("CONNECTION_MODE should be either internal or external")
		return AppCfg{}
	}
	if conf.ConnectionMode == "external" && conf.KubeconfigPath == "" {
		log.Fatalf("KUBECONFIG is required if CONNECTION_MODE is \"%s\"", conf.ConnectionMode)
		return AppCfg{}
	}

	return conf
}
