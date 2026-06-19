package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url-shortener/internal/app"
	"url-shortener/internal/config"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	app, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := app.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}
