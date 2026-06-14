package bootstrap

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/fx"

	"github.com/cubelitblade/community-v2/backend/services/post/internal/shared"
)

func provideDB(cfg shared.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.URI)
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
