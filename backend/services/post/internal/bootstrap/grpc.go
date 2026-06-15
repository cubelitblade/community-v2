package bootstrap

import (
	"context"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authzv1 "github.com/cubelitblade/community-v2/backend/gen/go/authz/v1"
)

func provideAuthzConn(lc fx.Lifecycle) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return conn.Close()
		},
	})

	return conn, nil
}

func provideAuthorizer(conn *grpc.ClientConn) authzv1.AuthzServiceClient {
	return authzv1.NewAuthzServiceClient(conn)
}

func GRPCModule() fx.Option {
	return fx.Options(
		fx.Provide(provideAuthzConn),
		fx.Provide(provideAuthorizer),
	)
}
