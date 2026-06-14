package bootstrap

import (
	"log/slog"

	"go.uber.org/fx"

	post "github.com/cubelitblade/community-v2/backend/services/post/internal/post/setup"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(LoadConfig),

		DBModule(),
		SlogModule(),

		post.Module(),

		fx.Invoke(func(logger *slog.Logger) {
			logger.Info("abc")
		}),
	)
}
