package touringparty

import (
	"context"

	"github.com/pranayhere/touringparty/internal/config"
	"github.com/pranayhere/touringparty/pkg/logger"
	"github.com/pranayhere/touringparty/service"
	"github.com/pranayhere/touringparty/service/web"
)

type Component struct {
	name    string
	service service.Service
}

func NewComponent(component string, cfg config.Config) (*Component, error) {
	var svc service.Service
	var err error

	switch component {
	case Web:
		svc, err = web.NewService(cfg)
	}

	if err != nil {
		return nil, err
	}

	return &Component{
		name:    component,
		service: svc,
	}, nil
}

// Run a touringparty component
func (c *Component) Run(ctx context.Context) error {
	logger.Ctx(ctx).Infow("starting touring party component", "name", c.name)
	return c.service.Start(ctx)
}
