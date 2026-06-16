package bootstrap

import (
	"context"
	"database/sql"
	"fmt"

	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"
	"go.uber.org/fx"

	"github.com/XSAM/otelsql"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/shared"
)

func provideDB(cfg shared.DatabaseConfig) (*sql.DB, error) {
	db, err := otelsql.Open("pgx", cfg.URI, otelsql.WithAttributes(semconv.DBSystemNamePostgreSQL))
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	return db, nil
}

func DBModule() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(provideDB,
				fx.OnStart(func(ctx context.Context, db *sql.DB) error { return db.PingContext(ctx) }),
				fx.OnStop(func(db *sql.DB) error { return db.Close() }),
			),
		),
	)
}
