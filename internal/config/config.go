package config

import (
	"github.com/pranayhere/touringparty/pkg/monitoring/sentry"
	"github.com/pranayhere/touringparty/pkg/tracing"
)

// Config is application config
type Config struct {
	App     App            `yaml:"app"`
	Server  Server         `yaml:"server"`
	Tracing tracing.Config `yaml:"tracing"`
	Sentry  sentry.Config  `yaml:"sentry"`
}

// App contains application-specific config values
type App struct {
	Env             string `yaml:"env"`
	ServiceName     string `yaml:"serviceName"`
	ShutdownTimeout int    `yaml:"shutdownTimeout"`
	ShoutdownDelay  int    `yaml:"shutdownDelay"`
	GitCommitHash   string `yaml:"gitCommitHash"`
}

type Server struct {
	Grpc struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
	Http struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
}
