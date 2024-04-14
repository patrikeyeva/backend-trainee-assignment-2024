package main

import (
	h "avito-banner/internal/handler"
	m "avito-banner/internal/middleware"
	"avito-banner/internal/repository/cache"
	"avito-banner/internal/repository/db"
	"avito-banner/internal/service"
	"context"
	"log"
	"net"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	//TODO: вынести в переменные окружения
	poolConfig, errPool := pgxpool.ParseConfig("postgres://postgres:postgres@localhost:5433/postgres")
	if errPool != nil {
		log.Fatal(errPool)
	}
	connDB, errDB := pgxpool.NewWithConfig(ctx, poolConfig)
	if errDB != nil {
		log.Fatal(errDB)
	}

	errPing := connDB.Ping(ctx)
	if errPing != nil {
		log.Fatal(errPing)
	}
	repo := db.NewRepository(connDB)
	cache := cache.NewCache()
	go cache.BackgroundCleaning(ctx)
	service := service.NewService(repo, cache)

	//TODO: вынести в переменные окружения
	conn, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	http.Handle("GET /user_banner", m.LoggingMiddleware(h.GetUserBanner(service)))
	http.Handle("GET /banner", m.LoggingMiddleware(h.GetBanner(service)))
	http.Handle("POST /banner", m.LoggingMiddleware(h.PostBanner(service)))
	http.Handle("PATCH /banner/{id}", m.LoggingMiddleware(h.PatchBannerId(service)))
	http.Handle("DELETE /banner/{id}", m.LoggingMiddleware(h.DeleteBannerId(service)))

	log.Println("Starting server")

	if errServ := http.Serve(conn, nil); err != nil {
		log.Fatal(errServ)
	}
}
