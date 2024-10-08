package config

import (
	"cmp"
	"os"

	"github.com/danielmesquitta/tasks-api/internal/pkg/validator"
	"github.com/spf13/viper"
)

type Environment string

const (
	DevelopmentEnv Environment = "development"
	TestEnv        Environment = "test"
	ProductionEnv  Environment = "production"
)

type Env struct {
	val validator.Validator

	Environment          Environment `mapstructure:"ENVIRONMENT"`
	Port                 string      `mapstructure:"PORT"`
	DBConnection         string      `mapstructure:"DB_CONNECTION"         validate:"required"`
	CipherSecretKey      string      `mapstructure:"CIPHER_SECRET_KEY"     validate:"required,min=32,max=32"`
	InitializationVector string      `mapstructure:"INITIALIZATION_VECTOR" validate:"required,min=16,max=16"`
	JWTSecretKey         string      `mapstructure:"JWT_SECRET_KEY"        validate:"required"`
	BasicAuthUsername    string      `mapstructure:"BASIC_AUTH_USERNAME"   validate:"required"`
	BasicAuthPassword    string      `mapstructure:"BASIC_AUTH_PASSWORD"   validate:"required"`
}

func (e *Env) validate() error {
	if err := e.val.Validate(e); err != nil {
		return err
	}
	if e.Environment == "" {
		e.Environment = DevelopmentEnv
	}
	if e.Port == "" {
		e.Port = "8080"
	}
	return nil
}

func LoadEnv(val validator.Validator) *Env {
	env := &Env{
		val: val,
	}

	envFilepath := cmp.Or(os.Getenv("ENV_FILEPATH"), ".env")

	viper.SetConfigFile(envFilepath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&env); err != nil {
		panic(err)
	}

	if err := env.validate(); err != nil {
		panic(err)
	}

	return env
}
