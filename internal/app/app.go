package app

import (
	"context"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/ports/http"
	"github.com/Rasikrr/core/application"
	"github.com/Rasikrr/core/log"
)

type App struct {
	application.App
}

func InitApp(ctx context.Context) *App {
	app := &App{
		App: *application.NewApp(ctx),
	}
	for _, initFn := range []func(context.Context) error{
		app.initHTTP,
	} {
		if err := initFn(ctx); err != nil {
			log.Fatal(ctx, "app init", log.Err(err))
		}
	}
	return app
}

func (a *App) initHTTP(ctx context.Context) error {
	http.NewServer(a.HTTPServer())
	return nil
}
