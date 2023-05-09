package main

import (
	"context"
	"flag"

	"github.com/pranayhere/touringparty/cmd/api/touringparty"
	"github.com/pranayhere/touringparty/internal/app"
	"github.com/pranayhere/touringparty/internal/boot"
)

var (
	componentName *string
)

func init() {
	componentName = flag.String("component", touringparty.Web, "component to start")
}

func main() {
	ctx := boot.NewContext(context.Background())

	// parse the cmd input
	flag.Parse()

	// read the env
	env := app.GetEnv()

	// Init app dependencies
	touringparty.Init(ctx, env, *componentName)

	touringparty.Run(ctx)
}
