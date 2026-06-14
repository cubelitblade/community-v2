package setup

import (
	"database/sql"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"go.uber.org/fx"

	"github.com/cubelitblade/community-v2/backend/services/post/internal/post/storage/ent"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(provideEntClient),
	)
}

func provideEntClient(db *sql.DB) *ent.Client {
	return ent.NewClient(
		ent.Driver(entsql.OpenDB(dialect.Postgres, db)),
	)
}
