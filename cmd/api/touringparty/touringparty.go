package touringparty

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pranayhere/touringparty/internal/boot"
	"github.com/pranayhere/touringparty/internal/config"
	configreader "github.com/pranayhere/touringparty/pkg/config"
	"github.com/pranayhere/touringparty/pkg/logger"
)

const (
	// Web component to which exposes APIs
	Web = "web"
	// Worker component which fires webhooks to subscribers
	Worker = "worker"
	// OpenAPIServer to server swagger docs
	OpenAPIServer = "openapi-server"
)

var validComponents = []string{Web, Worker, OpenAPIServer}
var component *Component

// isValidComponent validates if the input component is a valid metro component
// validComponents : web, worker
func isValidComponent(component string) bool {
	for _, s := range validComponents {
		if s == component {
			return true
		}
	}
	return false
}

// Init initializes all modules (logger, tracing, config, metro component)
func Init(_ context.Context, env string, componentName string) {
	// componentName validation
	ok := isValidComponent(componentName)
	if !ok {
		log.Fatalf("invalid componentName given as input: [%v]", componentName)
	}
	log.Printf("Setting up touringparty component: [%v] in env: [%v]", componentName, env)

	// read the componentName config for env
	var appConfig config.Config
	err := configreader.NewDefaultConfig().Load(env, &appConfig)
	if err != nil {
		log.Fatal(err)
	}

	if !ok {
		log.Fatalf("%v config missing", componentName)
	}

	err = boot.InitMonitoring(env, appConfig.App, appConfig.Sentry, appConfig.Tracing)

	if err != nil {
		log.Fatalf("error in setting up monitoring : %v", err)
	}

	// Init the requested componentName
	component, err = NewComponent(componentName, appConfig)
	if err != nil {
		log.Fatalf("error in creating touringparty component : %v", err)
	}
}

func Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	defer func() {
		err := boot.Close()

		if err != nil {
			log.Fatalf("error closing tracer: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	defer func() {
		signal.Stop(sigCh)
		cancel()
	}()

	go func() {
		sig := <-sigCh
		logger.Ctx(ctx).Infow("received a signal, stopping touring party", "signal", sig)
		errChan := component.service.GetErrorChannel()

		if errChan != nil {
			errChan <- fmt.Errorf("OS Signal error")
		}
		cancel()
	}()

	err := component.Run(ctx)
	if err != nil {
		logger.Ctx(ctx).Fatalw("component exited with error", "msg", err.Error())
	}

	logger.Ctx(ctx).Infow("stopped touringparty")
}
