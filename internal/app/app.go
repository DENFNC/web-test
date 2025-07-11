package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/DENFNC/web-test/config"
	"github.com/DENFNC/web-test/internal/infra/psql"
	"github.com/DENFNC/web-test/internal/infra/psql/repository"
	"github.com/DENFNC/web-test/internal/service"
	handler "github.com/DENFNC/web-test/internal/transport/http"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	*slog.Logger
	*http.ServeMux
	Addr string
}

func NewApp(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	mux := http.NewServeMux()

	db, err := initDatabase(ctx, cfg)
	if err != nil {
		log.Error(
			"Не удалось подключить к базе",
			slog.String("err", err.Error()),
		)
		os.Exit(1)
	}

	authRepo := repository.NewAuthRepository(log, db)
	authService := service.NewAuthService(log, authRepo)

	docRepo := repository.NewDocumentRepository(log, db)
	docService := service.NewDocumentService(docRepo, authRepo)

	handler.NewAuthHandler(log, mux, authService)
	handler.NewDocumentHandler(log, mux, docService)

	return &App{
		Logger:   log,
		ServeMux: mux,
		Addr:     cfg.AppConfig.URL,
	}
}

func (app *App) MustStart() {
	const op = "app.Start"

	log := app.Logger.With("op", op)

	log.Info(
		"Starting the server",
		slog.String("addr", app.Addr),
	)

	if err := http.ListenAndServe(app.Addr, app.ServeMux); err != nil {
		log.Error(
			"Failed to start tcp server",
			slog.String("err", err.Error()),
		)
		panic(err)
	}
}

func initDatabase(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	db, err := psql.NewDatabase(ctx, cfg.DBConfig.URL,
		psql.WithMaxConns(cfg.DBConfig.MaxConns),
		psql.WithMinConns(cfg.DBConfig.MinConns),
		psql.WithMaxConnIdleTime(cfg.DBConfig.MaxConnIdleTime),
		psql.WithMaxConnLifetime(cfg.DBConfig.MaxConnLifeTime),
		psql.WithHealthCheckPeriod(cfg.DBConfig.HealthCheckPeriod),
	)
	if err != nil {
		return nil, err
	}

	return db.Pool, nil
}
