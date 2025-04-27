package jwt

import (
	"log"

	"github.com/caarlos0/env"
)

type JwtConfig struct {
	JwtLifespanMinutes string `env:"JWT_LIFESPAN_MINUTES" envDefault:"10"`
	JwtSecret          string `env:"JWT_SECRET" envDefault:"1"`
}

func NewJwtConfig() (j *JwtConfig) {
	j = &JwtConfig{}
	err := env.Parse(j)
	if err != nil {
		log.Fatalf("Jwt config initialization error: %v", err)
	}
	return j
}
