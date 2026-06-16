package bootstrap

import (
	"database/sql"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"go.uber.org/fx"

	"github.com/cubelitblade/community-v2/backend/services/post/internal/adapter/driven/authz"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/adapter/driven/misc"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/adapter/driven/persistence"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/adapter/driven/persistence/ent"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/adapter/driving/grpc"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/application"
)

func PostModule() fx.Option {
	return fx.Options(
		fx.Provide(provideEntClient),

		fx.Provide(persistence.NewEntTransactor),
		fx.Provide(persistence.NewPostSaver),

		fx.Provide(authz.NewAuthorizer),

		fx.Provide(misc.NewIDGenerator),

		fx.Provide(grpc.NewPostServiceServer),

		fx.Provide(application.NewPublishPost),
	)
}

func provideEntClient(db *sql.DB) *ent.Client {
	return ent.NewClient(
		ent.Driver(entsql.OpenDB(dialect.Postgres, db)),
	)
}
