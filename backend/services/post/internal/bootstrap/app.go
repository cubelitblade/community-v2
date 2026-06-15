package bootstrap

import (
	"context"

	"github.com/cubelitblade/community-v2/backend/services/post/internal/application"
	"go.uber.org/fx"
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
