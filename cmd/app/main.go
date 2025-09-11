package main

import (
	"context"

	"github.com/Rasikrr/core/log"

	app2 "github.com/Rasikrr/bugsy_backend_monolith/internal/app"
)

func main() {
	ctx := context.Background()
	app := app2.InitApp(ctx)
	if err := app.Start(ctx); err != nil {
		log.Fatal(ctx, "app start", log.Err(err))
	}
}
