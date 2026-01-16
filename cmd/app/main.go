package main

import (
	"context"

	"github.com/Rasikrr/core/log"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/app"
)

func main() {
	ctx := context.Background()
	app := app.InitApp(ctx)
	if err := app.Start(ctx); err != nil {
		log.Fatal(ctx, "app start", log.Err(err))
	}
}
