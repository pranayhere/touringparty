package main

import (
	"context"

	"github.com/pranayhere/touringparty/cmd/api/touringparty"
	"github.com/pranayhere/touringparty/internal/boot"
)

func main() {
	ctx := boot.NewContext(context.Background())

	// env := app.GetEnv()

	touringparty.Run(ctx)
}
