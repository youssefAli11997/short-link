package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
)

func main() {
	ctx := context.Background()

	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	baseURL := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if baseURL == "" {
		log.Fatal("BASE_URL is required")
	}

	if port == "" {
		log.Fatal("PORT is required")
	}

	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	repository := repository.NewPostgresURLRepository(db)
	service := service.NewURLService(repository, baseURL)
	handler := handler.NewURLHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /encode", handler.Encode)
	mux.HandleFunc("POST /decode", handler.Decode)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: mux,
	}

	go func() {
		log.Println("server listening on :8080")

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	server.Shutdown(ctx)
}
