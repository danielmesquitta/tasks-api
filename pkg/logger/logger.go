package logger

import (
	"github.com/danielmesquitta/tasks-api/internal/config"
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger(env *config.Env) *Logger {
	var zapConfig zap.Config
	if env.Environment == config.TaskionEnv {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}

	zapConfig.DisableStacktrace = true

	return &Logger{
		SugaredLogger: zap.Must(zapConfig.Build()).Sugar(),
	}
}
