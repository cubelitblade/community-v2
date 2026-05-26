package bootstrap

import (
	"net/http"

	"github.com/cubelitblade/community-v2/backend/pkg/platform/logger"
	"github.com/cubelitblade/community-v2/backend/pkg/platform/server"
	authzSetup "github.com/cubelitblade/community-v2/backend/services/authz/internal/authz/setup"
	eventsSetup "github.com/cubelitblade/community-v2/backend/services/authz/internal/events/setup"
	healthSetup "github.com/cubelitblade/community-v2/backend/services/authz/internal/health/setup"
	"go.uber.org/fx"
)

// NewApp creates the fx application with all providers and lifecycle hooks.
func NewApp() *fx.App {
	return fx.New(
		fx.Provide(ProvideConfig),
		fx.Provide(logger.NewLogger),
		fx.Provide(server.NewGinEngine),
		fx.Provide(server.ProvideHTTPServer),

		authzSetup.Module(),
		healthSetup.Module(),
		eventsSetup.Module(),

		fx.Invoke(registerAPIRoutes),
		fx.Invoke(registerHealthRoutes),
		fx.Invoke(func(*http.Server) {}),
	)
}
