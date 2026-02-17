package main

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/app"
	"github.com/Rasikrr/core/log"
)

func main() {
	ctx := context.Background()

	a := app.InitApp(ctx)
	if err := a.Start(ctx); err != nil {
		log.Fatal(ctx, "start app error", log.Err(err))
	}
}
