// Package setup provides constructors that wire together the events subsystem.
package setup

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cubelitblade/community-v2/backend/pkg/events"
	authzEvents "github.com/cubelitblade/community-v2/backend/services/authz/internal/events"
	gosdk "github.com/openfga/go-sdk/client"
	"go.uber.org/fx"
)

// Module returns the fx providers for the events module.
func Module() fx.Option {
	return fx.Module(
		"events",
		fx.Provide(authzEvents.NewHandler),
		fx.Provide(NewWriter),
		fx.Provide(NewConsumer),
	)
}

// NewWriter adapts *gosdk.OpenFgaClient to the internal FGA WriteClient interface.
func NewWriter(client *gosdk.OpenFgaClient, logger *slog.Logger) *authzEvents.Writer {
	return authzEvents.NewWriter(client, logger)
}

// NewConsumer creates the events Consumer and registers lifecycle hooks.
func NewConsumer(cfg *events.Config, handler *authzEvents.Handler, logger *slog.Logger, lc fx.Lifecycle) (*events.Consumer, error) {
	consumer := events.NewConsumer(cfg, handler,5, logger)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := consumer.Start(ctx); err != nil {
				return fmt.Errorf("start consumer: %w", err)
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			consumer.Stop()
			return nil
		},
	})

	return consumer, nil
}
