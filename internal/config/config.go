package config

import (
	"github.com/pranayhere/touringparty/pkg/monitoring/sentry"
	"github.com/pranayhere/touringparty/pkg/tracing"
)

// Config is application config
type Config struct {
	App     App
	Tracing tracing.Config
	Sentry  sentry.Config
}

// App contains application-specific config values
type App struct {
	Env             string
	ServiceName     string
	ShutdownTimeout int
	ShoutdownDelay  int
	GitCommitHash   string
}
