package bootstrap

import (
	"context"

	"go.uber.org/fx"

	"github.com/cubelitblade/community-v2/backend/services/post/internal/application"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(LoadConfig),

		DBModule(),
		SlogModule(),
		GRPCModule(),
		SnowflakeModule(),

		PostModule(),

		fx.Invoke(func(p *application.PublishPost) {
			p.Publish(context.Background(), 1, nil, "content")
		}),
	)
}
