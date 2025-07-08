package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/DENFNC/web-test/config"
	"github.com/DENFNC/web-test/internal/repository"
	"github.com/DENFNC/web-test/internal/service"
	"github.com/DENFNC/web-test/internal/transport"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	// --- PostgreSQL ---
	pgDsn := "postgres://" + cfg.DBUser + ":" + cfg.DBPass + "@" + cfg.DBHost + ":" + cfg.DBPort + "/" + cfg.DBName + "?sslmode=disable"
	pgpool, err := pgxpool.New(context.Background(), pgDsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	if err := pgpool.Ping(context.Background()); err != nil {
		log.Fatalf("PostgreSQL не отвечает: %v", err)
	}

	// --- Redis ---
	redisCache := repository.NewRedisCache(cfg.RedisAddr)

	// --- Репозиторий ---
	repo := repository.NewPsqlRepo(pgpool)

	adminToken := os.Getenv("ADMIN_TOKEN")
	if adminToken == "" {
		adminToken = "superadmin123"
	}
	userSvc := service.NewUserService(repo, adminToken)
	docSvc := service.NewDocumentService(repo)
	server := transport.NewServerWithCache(repo, docSvc, userSvc, redisCache)
	log.Println("Сервер запущен на :8080 (PostgreSQL, Redis)")
	http.ListenAndServe(":8080", server)
}
