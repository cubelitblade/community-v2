// Package setup provides constructors that wire together authz domain services
// and HTTP transport.
package setup

import (
	"fmt"
	"log/slog"

	"github.com/cubelitblade/community-v2/backend/pkg/platform"
	"github.com/cubelitblade/community-v2/backend/services/authz/internal/authz"
	"github.com/cubelitblade/community-v2/backend/services/authz/internal/authz/transport"
	gosdk "github.com/openfga/go-sdk/client"
	"go.uber.org/fx"
)

// Module returns the fx providers for the authz module.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewFGAClient),
		fx.Provide(NewChecker),
		fx.Provide(
			fx.Annotate(NewHandler, fx.As(new(platform.HTTPMounter)), fx.ResultTags(`group:"mounter"`)),
		),
	)
}

// NewFGAClient creates the OpenFGA SDK client from config.
func NewFGAClient(cfg *authz.FGAConfig) (*gosdk.OpenFgaClient, error) {
	client, err := gosdk.NewSdkClient(&gosdk.ClientConfiguration{
		ApiUrl:               cfg.APIURL,
		StoreId:              cfg.StoreID,
		AuthorizationModelId: cfg.ModelID,
	})
	if err != nil {
		return nil, fmt.Errorf("create openfga client: %w", err)
	}

	return client, nil
}

// NewChecker creates the authz Checker backed by the OpenFGA SDK client.
func NewChecker(client *gosdk.OpenFgaClient) *authz.Checker {
	return authz.NewChecker(client)
}

// NewHandler creates the authz HTTP handler with all dependencies wired.
func NewHandler(checker *authz.Checker, logger *slog.Logger) *transport.Handler {
	return transport.NewHandler(transport.Deps{
		Checker: checker,
		Logger:  logger,
	})
}
