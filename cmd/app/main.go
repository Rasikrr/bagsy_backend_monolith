package main

import (
	"context"

	"github.com/Rasikrr/core/log"

	app2 "github.com/Rasikrr/bugsy_backend_monolith/internal/app"
)

// @title Bagsy API
// @version 1.0
// @description API для приложения Bagsy
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@bagsy.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /

func main() {
	ctx := context.Background()
	app := app2.InitApp(ctx)
	if err := app.Start(ctx); err != nil {
		log.Fatal(ctx, "app start", log.Err(err))
	}
}
